package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gocql/gocql"
)

const (
	accessTokenLifetime  = 15 * time.Minute
	refreshTokenLifetime = 7 * 24 * time.Hour
	refreshCookieName    = "refresh_token"
	refreshCookiePath    = "/auth"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func generateRandomToken(size int) (string, error) {
	randomBytes := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		return "", fmt.Errorf("generate random token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func encodeJWT(claims jwt.MapClaims, lifetime time.Duration) (string, error) {
	tokenClaims := make(jwt.MapClaims, len(claims)+2)
	for key, value := range claims {
		tokenClaims[key] = value
	}
	jwtauth.SetExpiryIn(tokenClaims, lifetime)
	jwtauth.SetIssuedNow(tokenClaims)

	_, tokenString, err := TokenAuth.Encode(tokenClaims)
	if err != nil {
		return "", fmt.Errorf("encode JWT: %w", err)
	}
	if tokenString == "" {
		return "", errors.New("encode JWT: empty token")
	}
	return tokenString, nil
}

func issueAccessToken(claims jwt.MapClaims) (string, error) {
	accessClaims := make(jwt.MapClaims, len(claims)+1)
	for key, value := range claims {
		accessClaims[key] = value
	}
	accessClaims["token_use"] = "access"
	return encodeJWT(accessClaims, accessTokenLifetime)
}

func requireAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil || claims["token_use"] != "access" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func requireSecureCredentialTransport(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if externalRequestScheme(r) != "https" && !AppConfig.JWT.ALLOWINSECUREHTTP {
			http.Error(w, "HTTPS is required for password operations", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func refreshTokenHash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func storeRefreshToken(token, userName string) error {
	expiresAt := time.Now().Add(refreshTokenLifetime)
	ttl := int(refreshTokenLifetime.Seconds())
	tokenHash := refreshTokenHash(token)
	if err := Session.Query(
		`INSERT INTO auth_refresh_sessions_by_user (user_name, token_hash) VALUES (?, ?) USING TTL ?`,
		userName, tokenHash, ttl,
	).WithContext(context.Background()).Exec(); err != nil {
		return fmt.Errorf("index refresh session by user: %w", err)
	}
	if err := Session.Query(
		`INSERT INTO auth_refresh_sessions (token_hash, user_name, expires_at, consumed) VALUES (?, ?, ?, false) USING TTL ?`,
		tokenHash, userName, expiresAt, ttl,
	).WithContext(context.Background()).Exec(); err != nil {
		if cleanupErr := deleteRefreshTokenIndex(userName, tokenHash); cleanupErr != nil {
			return fmt.Errorf("store refresh session: %w (also failed to remove user index: %v)", err, cleanupErr)
		}
		return fmt.Errorf("store refresh session: %w", err)
	}
	return nil
}

func consumeRefreshToken(token string) (bool, error) {
	var previousConsumed bool
	return Session.Query(
		`UPDATE auth_refresh_sessions SET consumed = true WHERE token_hash = ? IF consumed = false`,
		refreshTokenHash(token),
	).WithContext(context.Background()).ScanCAS(&previousConsumed)
}

func deleteRefreshToken(token string) error {
	var userName string
	err := Session.Query(
		`SELECT user_name FROM auth_refresh_sessions WHERE token_hash = ?`,
		refreshTokenHash(token),
	).WithContext(context.Background()).Consistency(ConsistencyRead).Scan(&userName)
	if err != nil && !errors.Is(err, gocql.ErrNotFound) {
		return fmt.Errorf("read refresh session before deletion: %w", err)
	}
	if err := Session.Query(
		`DELETE FROM auth_refresh_sessions WHERE token_hash = ?`,
		refreshTokenHash(token),
	).WithContext(context.Background()).Exec(); err != nil {
		return fmt.Errorf("delete refresh session: %w", err)
	}
	if userName != "" {
		if err := deleteRefreshTokenIndex(userName, refreshTokenHash(token)); err != nil {
			return err
		}
	}
	return nil
}

func deleteRefreshTokenIndex(userName, tokenHash string) error {
	if err := Session.Query(
		`DELETE FROM auth_refresh_sessions_by_user WHERE user_name = ? AND token_hash = ?`,
		userName, tokenHash,
	).WithContext(context.Background()).Exec(); err != nil {
		return fmt.Errorf("delete refresh-session user index: %w", err)
	}
	return nil
}

func revokeUserRefreshTokens(userName string) error {
	iter := Session.Query(
		`SELECT token_hash FROM auth_refresh_sessions_by_user WHERE user_name = ?`,
		userName,
	).WithContext(context.Background()).Consistency(gocql.Quorum).Iter()

	tokenHashes := make([]string, 0)
	var tokenHash string
	for iter.Scan(&tokenHash) {
		tokenHashes = append(tokenHashes, tokenHash)
	}
	if err := iter.Close(); err != nil {
		return fmt.Errorf("list refresh sessions for user: %w", err)
	}

	for _, hash := range tokenHashes {
		var previousConsumed bool
		applied, err := Session.Query(
			`UPDATE auth_refresh_sessions SET consumed = true WHERE token_hash = ? IF consumed = false`,
			hash,
		).WithContext(context.Background()).ScanCAS(&previousConsumed)
		if err != nil {
			return fmt.Errorf("revoke refresh session %s: %w", hash, err)
		}
		if err := deleteRefreshTokenIndex(userName, hash); err != nil {
			return err
		}
		if !applied {
			continue
		}
	}
	return nil
}

func setRefreshCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     refreshCookiePath,
		Expires:  time.Now().Add(refreshTokenLifetime),
		MaxAge:   int(refreshTokenLifetime.Seconds()),
		HttpOnly: true,
		Secure:   externalRequestScheme(r) == "https",
		SameSite: http.SameSiteStrictMode,
	})
}

func clearRefreshCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   externalRequestScheme(r) == "https",
		SameSite: http.SameSiteStrictMode,
	})
}

func clearLegacyAccessCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   externalRequestScheme(r) == "https",
		SameSite: http.SameSiteLaxMode,
	})
}

func externalRequestScheme(r *http.Request) string {
	if r.TLS != nil || AppConfig.HTTP.TLS {
		return "https"
	}
	if AppConfig.HTTP.TRUSTPROXY && strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")), "https") {
		return "https"
	}
	return "http"
}

func authRequestHasValidOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	parsedOrigin, err := url.Parse(origin)
	if err != nil || parsedOrigin.Host == "" || parsedOrigin.User != nil ||
		parsedOrigin.Path != "" || parsedOrigin.RawQuery != "" || parsedOrigin.Fragment != "" {
		return false
	}
	return strings.EqualFold(parsedOrigin.Scheme, externalRequestScheme(r)) &&
		strings.EqualFold(parsedOrigin.Host, r.Host)
}

func loginWithRefreshToken(w http.ResponseWriter, r *http.Request) {
	if externalRequestScheme(r) != "https" && !AppConfig.JWT.ALLOWINSECUREHTTP {
		http.Error(w, "HTTPS is required for password operations", http.StatusBadRequest)
		return
	}
	if !authRequestHasValidOrigin(r) {
		http.Error(w, "invalid request origin", http.StatusForbidden)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid login request", http.StatusBadRequest)
		return
	}

	userName := r.FormValue("username")
	claims, err := authenticateUser(userName, r.FormValue("password"))
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken, err := issueAccessToken(claims)
	if err != nil {
		log.Printf("Failed to issue access token: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	refreshToken, err := generateRandomToken(32)
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := storeRefreshToken(refreshToken, userName); err != nil {
		log.Printf("Failed to store refresh token: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := recordSuccessfulLogin(userName); err != nil {
		if cleanupErr := deleteRefreshToken(refreshToken); cleanupErr != nil {
			log.Printf("Failed to remove refresh token after login audit failure: %v", cleanupErr)
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeTokenResponse(w, r, refreshToken, accessToken)
}

func refreshAccessToken(w http.ResponseWriter, r *http.Request) {
	if !authRequestHasValidOrigin(r) {
		http.Error(w, "invalid request origin", http.StatusForbidden)
		return
	}

	cookie, err := r.Cookie(refreshCookieName)
	if err != nil || cookie.Value == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var userName string
	var expiresAt time.Time
	var consumed bool
	err = Session.Query(
		`SELECT user_name, expires_at, consumed FROM auth_refresh_sessions WHERE token_hash = ?`,
		refreshTokenHash(cookie.Value),
	).WithContext(r.Context()).Consistency(ConsistencyRead).Scan(&userName, &expiresAt, &consumed)
	if errors.Is(err, gocql.ErrNotFound) || (err == nil && (consumed || !time.Now().Before(expiresAt))) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Printf("Failed to read refresh session: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	claims, active, err := activeUserClaims(userName)
	if err != nil {
		log.Printf("Failed to load refresh-session user: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !active {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken, err := issueAccessToken(claims)
	if err != nil {
		log.Printf("Failed to issue access token: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	newRefreshToken, err := generateRandomToken(32)
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := storeRefreshToken(newRefreshToken, userName); err != nil {
		log.Printf("Failed to store rotated refresh token: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	applied, err := consumeRefreshToken(cookie.Value)
	if err != nil {
		log.Printf("Failed to rotate refresh token: %v", err)
		if cleanupErr := deleteRefreshToken(newRefreshToken); cleanupErr != nil {
			log.Printf("Failed to remove unused refresh token: %v", cleanupErr)
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := deleteRefreshTokenIndex(userName, refreshTokenHash(cookie.Value)); err != nil {
		log.Printf("Failed to remove rotated refresh-token index: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !applied {
		if cleanupErr := deleteRefreshToken(newRefreshToken); cleanupErr != nil {
			log.Printf("Failed to remove unused refresh token: %v", cleanupErr)
		}
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	writeTokenResponse(w, r, newRefreshToken, accessToken)
}

func logoutWithRefreshToken(w http.ResponseWriter, r *http.Request) {
	if !authRequestHasValidOrigin(r) {
		http.Error(w, "invalid request origin", http.StatusForbidden)
		return
	}

	cookie, err := r.Cookie(refreshCookieName)
	clearRefreshCookie(w, r)
	clearLegacyAccessCookie(w, r)
	if err == nil && cookie.Value != "" {
		var userName string
		if lookupErr := Session.Query(
			`SELECT user_name FROM auth_refresh_sessions WHERE token_hash = ?`,
			refreshTokenHash(cookie.Value),
		).WithContext(r.Context()).Consistency(ConsistencyRead).Scan(&userName); lookupErr != nil && !errors.Is(lookupErr, gocql.ErrNotFound) {
			log.Printf("Failed to find refresh session for logout: %v", lookupErr)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if _, revokeErr := consumeRefreshToken(cookie.Value); revokeErr != nil {
			log.Printf("Failed to revoke refresh token: %v", revokeErr)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if userName != "" {
			if indexErr := deleteRefreshTokenIndex(userName, refreshTokenHash(cookie.Value)); indexErr != nil {
				log.Printf("Failed to remove refresh-token user index during logout: %v", indexErr)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func activeUserClaims(userName string) (jwt.MapClaims, bool, error) {
	var userID string
	var status int
	var role string
	var ownerID string
	err := Session.Query(
		`SELECT user_id, status, role, owner_id FROM user WHERE user_name = ?`,
		userName,
	).WithContext(context.Background()).Consistency(ConsistencyRead).Scan(&userID, &status, &role, &ownerID)
	if errors.Is(err, gocql.ErrNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if status != 1 {
		return nil, false, nil
	}
	return jwt.MapClaims{"username": userName, "userid": userID, "role": role, "ownerid": ownerID}, true, nil
}

func writeTokenResponse(w http.ResponseWriter, r *http.Request, refreshToken, accessToken string) {
	response, err := json.Marshal(tokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(accessTokenLifetime.Seconds()),
	})
	if err != nil {
		log.Printf("Failed to marshal token response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	setRefreshCookie(w, r, refreshToken)
	clearLegacyAccessCookie(w, r)
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Printf("Failed to write token response: %v", err)
	}
}

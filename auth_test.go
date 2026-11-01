package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth/v5"
)

func TestGenerateRandomToken(t *testing.T) {
	token, err := generateRandomToken(32)
	if err != nil {
		t.Fatalf("generateRandomToken() error = %v", err)
	}
	if len(token) != 43 {
		t.Fatalf("generateRandomToken() length = %d, want 43", len(token))
	}
	if strings.Contains(token, "=") {
		t.Fatal("generateRandomToken() returned padded base64")
	}

	anotherToken, err := generateRandomToken(32)
	if err != nil {
		t.Fatalf("generateRandomToken() second call error = %v", err)
	}
	if token == anotherToken {
		t.Fatal("generateRandomToken() returned the same token twice")
	}
}

func TestIssueAccessToken(t *testing.T) {
	previousTokenAuth := TokenAuth
	TokenAuth = jwtauth.New("HS256", []byte(strings.Repeat("s", 32)), nil)
	t.Cleanup(func() {
		TokenAuth = previousTokenAuth
	})

	claims := jwt.MapClaims{"username": "alice", "role": "admin"}
	tokenString, err := issueAccessToken(claims)
	if err != nil {
		t.Fatalf("issueAccessToken() error = %v", err)
	}
	if _, exists := claims["token_use"]; exists {
		t.Fatal("issueAccessToken() mutated the input claims")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(strings.Repeat("s", 32)), nil
	})
	if err != nil {
		t.Fatalf("parse access token: %v", err)
	}
	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("token claims type = %T, want jwt.MapClaims", token.Claims)
	}
	if tokenClaims["token_use"] != "access" || tokenClaims["username"] != "alice" {
		t.Fatalf("unexpected access token claims: %v", tokenClaims)
	}
	if !tokenClaims.VerifyExpiresAt(time.Now().Add(14*time.Minute).Unix(), true) {
		t.Fatal("access token should remain valid for at least 14 minutes")
	}
	if tokenClaims.VerifyExpiresAt(time.Now().Add(16*time.Minute).Unix(), true) {
		t.Fatal("access token should expire within 16 minutes")
	}
}

func TestAuthRequestHasValidOrigin(t *testing.T) {
	request := httptest.NewRequest("POST", "https://example.com/auth/refresh", nil)
	request.Header.Set("Origin", "https://example.com")
	if !authRequestHasValidOrigin(request) {
		t.Fatal("same-origin request should be accepted")
	}

	request.Header.Set("Origin", "https://attacker.example")
	if authRequestHasValidOrigin(request) {
		t.Fatal("cross-origin request should be rejected")
	}
}

func TestProtectedRoutesRequireBearerAccessToken(t *testing.T) {
	previousTokenAuth := TokenAuth
	TokenAuth = jwtauth.New("HS256", []byte(strings.Repeat("s", 32)), nil)
	t.Cleanup(func() {
		TokenAuth = previousTokenAuth
	})

	protected := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler := jwtauth.Verify(TokenAuth, jwtauth.TokenFromHeader)(
		jwtauth.Authenticator(TokenAuth)(requireAccessToken(protected)),
	)

	makeToken := func(tokenUse string) string {
		claims := jwt.MapClaims{"username": "alice"}
		if tokenUse != "" {
			claims["token_use"] = tokenUse
		}
		jwtauth.SetIssuedNow(claims)
		jwtauth.SetExpiryIn(claims, time.Minute)
		_, tokenString, err := TokenAuth.Encode(claims)
		if err != nil {
			t.Fatalf("encode test token: %v", err)
		}
		return tokenString
	}

	tests := []struct {
		name          string
		authorization string
		cookie        string
		wantStatus    int
	}{
		{
			name:          "access bearer token",
			authorization: "Bearer " + makeToken("access"),
			wantStatus:    http.StatusNoContent,
		},
		{
			name:          "legacy token in bearer header",
			authorization: "Bearer " + makeToken(""),
			wantStatus:    http.StatusUnauthorized,
		},
		{
			name:          "refresh token in bearer header",
			authorization: "Bearer " + makeToken("refresh"),
			wantStatus:    http.StatusUnauthorized,
		},
		{
			name:       "access token in legacy cookie",
			cookie:     "jwt=" + makeToken("access"),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/v1/userinfo", nil)
			if test.authorization != "" {
				request.Header.Set("Authorization", test.authorization)
			}
			if test.cookie != "" {
				request.Header.Set("Cookie", test.cookie)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
		})
	}
}

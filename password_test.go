package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPasswordHashAndVerification(t *testing.T) {
	password := "correct horse battery staple"
	storedHash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword() error = %v", err)
	}
	if !strings.HasPrefix(storedHash, "$argon2id$v=19$") {
		t.Fatalf("hashPassword() format = %q, want Argon2id PHC format", storedHash)
	}

	verified, needsUpgrade, err := verifyPassword(storedHash, password)
	if err != nil {
		t.Fatalf("verifyPassword() error = %v", err)
	}
	if !verified || needsUpgrade {
		t.Fatalf("verifyPassword() = (%v, %v), want (true, false)", verified, needsUpgrade)
	}

	verified, _, err = verifyPassword(storedHash, "incorrect password")
	if err != nil {
		t.Fatalf("verifyPassword() with incorrect password error = %v", err)
	}
	if verified {
		t.Fatal("verifyPassword() accepted an incorrect password")
	}
}

func TestLegacySHA256PasswordVerification(t *testing.T) {
	password := "legacy password"
	legacyHash := sha256.Sum256([]byte(password))
	storedHash := fmt.Sprintf("%x", legacyHash)

	verified, needsUpgrade, err := verifyPassword(storedHash, password)
	if err != nil {
		t.Fatalf("verifyPassword() error = %v", err)
	}
	if !verified || !needsUpgrade {
		t.Fatalf("verifyPassword() = (%v, %v), want (true, true)", verified, needsUpgrade)
	}
}

func TestMalformedArgon2HashIsReported(t *testing.T) {
	if _, _, err := verifyPassword("$argon2id$invalid", "password"); err == nil {
		t.Fatal("verifyPassword() accepted a malformed Argon2id hash")
	}
}

func TestValidateNewPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{name: "minimum length", password: "12345678"},
		{name: "too short", password: "1234567", wantErr: true},
		{name: "too long", password: strings.Repeat("a", maximumPasswordBytes+1), wantErr: true},
		{name: "unicode minimum length", password: "пароль123"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateNewPassword(test.password)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateNewPassword() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestRequireSecureCredentialTransport(t *testing.T) {
	previousConfig := AppConfig
	AppConfig.HTTP.TLS = false
	AppConfig.HTTP.TRUSTPROXY = false
	AppConfig.JWT.ALLOWINSECUREHTTP = false
	t.Cleanup(func() {
		AppConfig = previousConfig
	})

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler := requireSecureCredentialTransport(next)

	request := httptest.NewRequest(http.MethodPost, "http://example.com/auth", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("HTTP status = %d, want %d", response.Code, http.StatusBadRequest)
	}

	request = httptest.NewRequest(http.MethodPost, "http://example.com/auth", nil)
	request.Header.Set("X-Forwarded-Proto", "https")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("untrusted forwarded HTTPS status = %d, want %d", response.Code, http.StatusBadRequest)
	}

	AppConfig.HTTP.TRUSTPROXY = true
	request = httptest.NewRequest(http.MethodPost, "http://example.com/auth", nil)
	request.Header.Set("X-Forwarded-Proto", "https")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("trusted forwarded HTTPS status = %d, want %d", response.Code, http.StatusNoContent)
	}

	request = httptest.NewRequest(http.MethodPost, "http://example.com/auth", nil)
	request.Header.Set("X-Forwarded-Proto", "https, http")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("ambiguous forwarded HTTPS status = %d, want %d", response.Code, http.StatusBadRequest)
	}

	AppConfig.HTTP.TRUSTPROXY = false
	AppConfig.JWT.ALLOWINSECUREHTTP = true
	request = httptest.NewRequest(http.MethodPost, "http://example.com/auth", nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("explicit test HTTP status = %d, want %d", response.Code, http.StatusNoContent)
	}
}

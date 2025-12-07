package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Fry-Fr/chirpy/internal/auth"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	tests := []struct {
		name      string
		userID    uuid.UUID
		secret    string
		expiresIn time.Duration
		wantErr   bool
	}{
		{
			name:      "valid JWT creation",
			userID:    uuid.New(),
			secret:    "test-secret",
			expiresIn: time.Hour,
			wantErr:   false,
		},
		{
			name:      "short expiration",
			userID:    uuid.New(),
			secret:    "test-secret",
			expiresIn: time.Second,
			wantErr:   false,
		},
		{
			name:      "empty secret",
			userID:    uuid.New(),
			secret:    "",
			expiresIn: time.Hour,
			wantErr:   false, // JWT library may allow this
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := auth.MakeJWT(tt.userID, tt.secret, tt.expiresIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("MakeJWT() returned empty token")
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"

	validToken, err := auth.MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create valid token: %v", err)
	}

	expiredToken, err := auth.MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	differentSecretToken, err := auth.MakeJWT(userID, "different-secret", time.Hour)
	if err != nil {
		t.Fatalf("Failed to create token with different secret: %v", err)
	}

	tests := []struct {
		name       string
		token      string
		secret     string
		wantUserID uuid.UUID
		wantErr    bool
	}{
		{
			name:       "valid token",
			token:      validToken,
			secret:     secret,
			wantUserID: userID,
			wantErr:    false,
		},
		{
			name:       "expired token",
			token:      expiredToken,
			secret:     secret,
			wantUserID: uuid.Nil,
			wantErr:    true,
		},
		{
			name:       "wrong secret",
			token:      differentSecretToken,
			secret:     secret,
			wantUserID: uuid.Nil,
			wantErr:    true,
		},
		{
			name:       "malformed token",
			token:      "invalid.token.here",
			secret:     secret,
			wantUserID: uuid.Nil,
			wantErr:    true,
		},
		{
			name:       "empty token",
			token:      "",
			secret:     secret,
			wantUserID: uuid.Nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := auth.ValidateJWT(tt.token, tt.secret)
			if err != nil && !tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v name %v", err, tt.wantErr, tt.name)
				return
			}
			if !tt.wantErr && gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() userID = %v, want %v", gotUserID, tt.wantUserID)
			}
			fmt.Printf("Got userID: %v\n", gotUserID.String())
		})
	}
}

func TestJWTRoundTrip(t *testing.T) {
	userID := uuid.New()
	secret := "round-trip-secret"
	expiresIn := time.Minute

	token, err := auth.MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() failed: %v", err)
	}

	returnedUserID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() failed: %v", err)
	}

	if returnedUserID != userID {
		t.Errorf("Round trip failed: expected %v, got %v", userID, returnedUserID)
	}
}

func TestGetBearerToken(t *testing.T) {
	httpHeaders := http.Header{}
	httpHeaders.Set("Authorization", "Bearer test-token-123")

	token, err := auth.GetBearerToken(httpHeaders)
	if err != nil {
		t.Fatalf("GetBearerToken() failed: %v", err)
	}

	expectedToken := "test-token-123"
	if token != expectedToken {
		t.Errorf("GetBearerToken() = %v, want %v", token, expectedToken)
	}
}

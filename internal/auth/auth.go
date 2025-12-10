package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/Fry-Fr/chirpy/internal/config"
	"github.com/Fry-Fr/chirpy/internal/database"
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	params := &argon2id.Params{
		Memory:      128 * 1024,
		Iterations:  4,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  16,
		KeyLength:   32,
	}
	hashed_pw, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", err
	}
	return hashed_pw, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Issuer:    "chirpy",
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, err
		}
		return userID, nil
	} else {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", http.ErrNoCookie
	}

	var tokenType, token string
	_, err := fmt.Sscanf(authHeader, "%s %s", &tokenType, &token)
	if err != nil || tokenType != "Bearer" {
		return "", http.ErrNoCookie
	}

	return token, nil
}

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

func MakeRefreshToken(cfg *config.ApiConfig, userId uuid.UUID) (string, error) {
	key := make([]byte, 32)
	rand.Read(key)
	encodedStr := hex.EncodeToString(key)

	refresh_tkn := encodedStr
	expires_in_hours := 1440 // default 60 days
	create_refresh_token_params := database.CreateRefreshTokenParams{
		Token:     refresh_tkn,
		UserID:    userId,
		ExpiresAt: time.Now().Add(time.Duration(expires_in_hours) * time.Hour),
	}
	refreshToken, err := cfg.DB.CreateRefreshToken(context.Background(), create_refresh_token_params)
	if err != nil {
		return "", err
	}
	return refreshToken.Token, nil
}

func GetRefreshToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", http.ErrNoCookie
	}

	var tokenType, token string
	_, err := fmt.Sscanf(authHeader, "%s %s", &tokenType, &token)
	if err != nil || tokenType != "Bearer" {
		return "", http.ErrNoCookie
	}

	return token, nil
}

func GetPolkaKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", http.ErrNoCookie
	}

	var tokenType, token string
	_, err := fmt.Sscanf(authHeader, "%s %s", &tokenType, &token)
	if err != nil || tokenType != "ApiKey" {
		return "", http.ErrNoCookie
	}

	if token != os.Getenv("POLKA_KEY") {
		return "", http.ErrNoCookie
	}

	return token, nil
	// return os.Getenv("POLKA_KEY")
}

package config

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// GetJWTSecret returns the secret used to sign JWTs. It reads the JWT_SECRET
// environment variable and falls back to a default (change for production).
func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "replace_this_with_a_secure_secret"
	}
	return secret
}

// ValidateJWT parses and validates the given token string using the configured secret.
// On success it returns the token claims as jwt.MapClaims.
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// ensure signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(GetJWTSecret()), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

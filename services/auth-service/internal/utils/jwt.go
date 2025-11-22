package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenExp  = 15 * time.Minute
	refreshTokenExp = 7 * 24 * time.Hour
)

var jwtSecret []byte

func InitSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "anime-platform-dev-secret-12345" // ← для демо
	}
	jwtSecret = []byte(secret)
}

func SignAccessToken(username string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(accessTokenExp).Unix(),
		"aud": "anime-platform",
	}).SignedString(jwtSecret)
}

func SignRefreshToken(username string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(refreshTokenExp).Unix(),
		"aud": "anime-platform",
	}).SignedString(jwtSecret)
}

func GetRefreshTokenExp() time.Duration {
	return refreshTokenExp
}

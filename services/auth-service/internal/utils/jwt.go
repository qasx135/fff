package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
	jwtSecret       []byte
}

func (m *JwtManager) InitSecret(secret string, accessTokenExpipartion time.Duration, refreshTokenExpiration time.Duration) {
	m.jwtSecret = []byte(secret)
	m.accessTokenExp = accessTokenExpipartion
	m.refreshTokenExp = refreshTokenExpiration
}

func (m *JwtManager) SignAccessToken(username string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(m.accessTokenExp).Unix(),
		"aud": "anime-platform",
	}).SignedString(m.jwtSecret)
}

func (m *JwtManager) SignRefreshToken(username string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(m.refreshTokenExp).Unix(),
		"aud": "anime-platform",
	}).SignedString(m.jwtSecret)
}

func (m *JwtManager) GetRefreshTokenExp() time.Duration {
	return m.refreshTokenExp
}

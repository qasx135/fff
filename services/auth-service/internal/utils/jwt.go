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
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   username,
		Audience:  []string{"anime-platform"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenExp)),
	}).SignedString(m.jwtSecret)
}

func (m *JwtManager) SignRefreshToken(username string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshTokenExp)),
		Audience:  []string{"anime-platform"},
	}).SignedString(m.jwtSecret)
}

func (m *JwtManager) GetRefreshTokenExp() time.Duration {
	return m.refreshTokenExp
}

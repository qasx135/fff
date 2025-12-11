package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	jwtSecret []byte
}

func (m *JwtManager) InitSecret(secret string) {
	m.jwtSecret = []byte(secret)
}

// Keyfunc callback
func (m *JwtManager) getJwtSecret(token *jwt.Token) (any, error) {
	// Ensure the signing method is HMAC (HS256 in this case)
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}
	return m.jwtSecret, nil
}

func (m *JwtManager) ValidateJWT(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{} // Initialize an empty claims struct

	// Pass the address of the custom claims struct to ParseWithClaims
	token, err := jwt.ParseWithClaims(tokenString, claims, m.getJwtSecret,
		jwt.WithValidMethods([]string{"HS256"}), // Explicitly define valid signing methods
		jwt.WithAudience("anime-platform"),      // Let the library handle audience validation
	)

	if err != nil {
		return "", err
	}

	rc, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", errors.New("invalid claims type")
	}

	// The claims are already populated in the 'claims' variable
	return rc.Subject, nil
}

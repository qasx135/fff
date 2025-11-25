package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// Define a custom claims struct
type MyClaims struct {
	Username string `json:"sub"`
	Audience string `json:"aud"`
	jwt.RegisteredClaims
}

type JwtManager struct {
	jwtSecret       []byte
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
	claims := &MyClaims{} // Initialize an empty claims struct

	// Pass the address of the custom claims struct to ParseWithClaims
	token, err := jwt.ParseWithClaims(tokenString, claims, m.getJwtSecret,
		jwt.WithValidMethods([]string{"HS256"}), // Explicitly define valid signing methods
		jwt.WithAudience("anime-platform"), // Let the library handle audience validation
	)

	if err != nil {
		return "", err
	}

	if token.Valid {
		// The claims are already populated in the 'claims' variable
		return claims.Username, nil
	}

	return "", errors.New("invalid token")
}

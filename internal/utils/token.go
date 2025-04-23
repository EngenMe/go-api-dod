package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenManager provides methods for creating and validating JWT tokens
type TokenManager struct {
	Secret    []byte
	Issuer    string
	ExpiresIn time.Duration
}

// NewTokenManager creates a new TokenManager
func NewTokenManager(
	secret, issuer string,
	expiresIn time.Duration,
) *TokenManager {
	return &TokenManager{
		Secret:    []byte(secret),
		Issuer:    issuer,
		ExpiresIn: expiresIn,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// Generate generates a JWT token for a user
func (m *TokenManager) Generate(userID uuid.UUID, email string) (
	string,
	error,
) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ExpiresIn)),
			Issuer:    m.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.Secret)
}

// Validate validates a JWT token
func (m *TokenManager) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return m.Secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

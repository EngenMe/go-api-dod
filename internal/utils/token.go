package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenManager provides methods for creating and validating JWT tokens
type TokenManager struct {
	Secret                []byte
	Issuer                string
	AccessTokenExpiresIn  time.Duration
	RefreshTokenExpiresIn time.Duration
}

// NewTokenManager creates a new TokenManager
func NewTokenManager(
	secret, issuer string,
	accessTokenExpiresIn, refreshTokenExpiresIn time.Duration,
) *TokenManager {
	return &TokenManager{
		Secret:                []byte(secret),
		Issuer:                issuer,
		AccessTokenExpiresIn:  accessTokenExpiresIn,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,
	}
}

// TokenType represents the type of token
type TokenType string

const (
	// AccessToken is a short-lived token used for API authentication
	AccessToken TokenType = "access"
	// RefreshToken is a long-lived token used to get new access tokens
	RefreshToken TokenType = "refresh"
)

// Claims represents JWT claims
type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a short-lived access token for a user
func (m *TokenManager) GenerateAccessToken(userID uuid.UUID, email string) (
	string,
	error,
) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Email:     email,
		TokenType: AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.AccessTokenExpiresIn)),
			Issuer:    m.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.Secret)
}

// GenerateRefreshToken generates a long-lived refresh token for a user
func (m *TokenManager) GenerateRefreshToken(userID uuid.UUID, email string) (
	string,
	error,
) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		Email:     email,
		TokenType: RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.RefreshTokenExpiresIn)),
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

// ValidateAccessToken validates an access token
func (m *TokenManager) ValidateAccessToken(tokenString string) (
	*Claims,
	error,
) {
	claims, err := m.Validate(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != AccessToken {
		return nil, errors.New("token is not an access token")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (m *TokenManager) ValidateRefreshToken(tokenString string) (
	*Claims,
	error,
) {
	claims, err := m.Validate(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != RefreshToken {
		return nil, errors.New("token is not a refresh token")
	}

	return claims, nil
}

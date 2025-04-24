package store

import (
	"time"

	"github.com/EngenMe/go-api-dod/internal/data/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshTokenStore provides methods to interact with the refresh_tokens table
type RefreshTokenStore struct {
	DB *gorm.DB
}

// NewRefreshTokenStore creates a new RefreshTokenStore
func NewRefreshTokenStore(db *gorm.DB) *RefreshTokenStore {
	return &RefreshTokenStore{
		DB: db,
	}
}

// Create creates a new refresh token
func (s *RefreshTokenStore) Create(refreshToken *models.RefreshToken) error {
	if refreshToken.ID == uuid.Nil {
		refreshToken.ID = uuid.New()
	}
	refreshToken.CreatedAt = time.Now()
	refreshToken.UpdatedAt = time.Now()

	query := `
        INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	result := s.DB.Exec(
		query,
		refreshToken.ID,
		refreshToken.UserID,
		refreshToken.Token,
		refreshToken.ExpiresAt,
		refreshToken.CreatedAt,
		refreshToken.UpdatedAt,
	)
	return result.Error
}

// GetByToken retrieves a refresh token by its token string
func (s *RefreshTokenStore) GetByToken(token string) (
	*models.RefreshToken,
	error,
) {
	var refreshToken models.RefreshToken
	query := `
        SELECT id, user_id, token, expires_at, created_at, updated_at, revoked_at
        FROM refresh_tokens
        WHERE token = $1
    `
	result := s.DB.Raw(query, token).Scan(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &refreshToken, nil
}

// GetByUserID retrieves all refresh tokens for a user
func (s *RefreshTokenStore) GetByUserID(userID uuid.UUID) (
	[]models.RefreshToken,
	error,
) {
	var refreshTokens []models.RefreshToken
	query := `
        SELECT id, user_id, token, expires_at, created_at, updated_at, revoked_at
        FROM refresh_tokens
        WHERE user_id = $1 AND revoked_at IS NULL
        ORDER BY created_at DESC
    `
	result := s.DB.Raw(query, userID).Scan(&refreshTokens)
	if result.Error != nil {
		return nil, result.Error
	}

	return refreshTokens, nil
}

// Revoke revokes a refresh token
func (s *RefreshTokenStore) Revoke(id uuid.UUID) error {
	now := time.Now()
	query := `
        UPDATE refresh_tokens
        SET revoked_at = $1, updated_at = $2
        WHERE id = $3
    `
	result := s.DB.Exec(query, now, now, id)
	return result.Error
}

// RevokeAllForUser revokes all refresh tokens for a user
func (s *RefreshTokenStore) RevokeAllForUser(userID uuid.UUID) error {
	now := time.Now()
	query := `
        UPDATE refresh_tokens
        SET revoked_at = $1, updated_at = $2
        WHERE user_id = $3 AND revoked_at IS NULL
    `
	result := s.DB.Exec(query, now, now, userID)
	return result.Error
}

// DeleteExpired deletes all expired refresh tokens
func (s *RefreshTokenStore) DeleteExpired() error {
	now := time.Now()
	query := `
        DELETE FROM refresh_tokens
        WHERE expires_at < $1
    `
	result := s.DB.Exec(query, now)
	return result.Error
}

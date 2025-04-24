package models

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a refresh token in the system
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	Token     string    `gorm:"type:varchar(512);uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	RevokedAt *time.Time
}

func (rt *RefreshToken) IsValid() bool {
	// Check if the token has expired
	if rt.ExpiresAt.Before(time.Now()) {
		return false
	}

	return true
}

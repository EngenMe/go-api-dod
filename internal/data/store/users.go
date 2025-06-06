package store

import (
	"time"

	"github.com/EngenMe/go-api-dod/internal/data/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserStore provides methods to interact with the user's table
type UserStore struct {
	DB *gorm.DB
}

// NewUserStore creates a new UserStore
func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		DB: db,
	}
}

// Create creates a new user
func (s *UserStore) Create(user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
        INSERT INTO users (id, email, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `
	result := s.DB.Exec(
		query,
		user.ID,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return result.Error
}

// GetByID retrieves a user by ID
func (s *UserStore) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, email, password, created_at, updated_at, deleted_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `
	result := s.DB.Raw(query, id).Scan(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (s *UserStore) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, email, password, created_at, updated_at, deleted_at
        FROM users
        WHERE email = $1 AND deleted_at IS NULL
    `
	result := s.DB.Raw(query, email).Scan(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &user, nil
}

// List retrieves all users
func (s *UserStore) List(limit, offset int) ([]models.User, error) {
	var users []models.User
	query := `
        SELECT id, email, password, created_at, updated_at, deleted_at
        FROM users
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
	result := s.DB.Raw(query, limit, offset).Scan(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// Update updates a user
func (s *UserStore) Update(user *models.User) error {
	user.UpdatedAt = time.Now()
	query := `
        UPDATE users
        SET email = $1,
            password = $2,
            updated_at = $3
        WHERE id = $4 AND deleted_at IS NULL
    `
	result := s.DB.Exec(
		query,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.ID,
	)

	return result.Error
}

// Delete deletes a user
func (s *UserStore) Delete(id uuid.UUID) error {
	query := `
        UPDATE users
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at IS NULL
    `
	result := s.DB.Exec(query, time.Now(), id)
	return result.Error
}

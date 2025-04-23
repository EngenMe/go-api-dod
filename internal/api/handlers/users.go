package handlers

import (
	"net/http"
	"strconv"

	"github.com/EngenMe/go-api-dod/internal/data/models"
	"github.com/EngenMe/go-api-dod/internal/data/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler provides handlers for user-related endpoints
type UserHandler struct {
	UserStore *store.UserStore
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userStore *store.UserStore) *UserHandler {
	return &UserHandler{
		UserStore: userStore,
	}
}

// CreateUser handles the creation of a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	// Parse request body
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	existingUser, err := h.UserStore.GetByEmail(req.Email)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to check user existence"},
		)
		return
	}
	if existingUser != nil {
		c.JSON(
			http.StatusConflict,
			gin.H{"error": "User with this email already exists"},
		)
		return
	}

	// Create user
	user := models.User{
		Email:    req.Email,
		Password: req.Password, // This will be hashed in the auth handler
	}

	if err := h.UserStore.Create(&user); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create user"},
		)
		return
	}

	// Return created user
	c.JSON(
		http.StatusCreated, gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	)
}

// GetUser handles retrieving a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	// Parse user ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user
	user, err := h.UserStore.GetByID(id)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get user"},
		)
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return user
	c.JSON(
		http.StatusOK, gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	)
}

// ListUsers handles retrieving all users
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get users
	users, err := h.UserStore.List(limit, offset)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to list users"},
		)
		return
	}

	// Map users to response
	var response []gin.H
	for _, user := range users {
		response = append(
			response, gin.H{
				"id":         user.ID,
				"email":      user.Email,
				"created_at": user.CreatedAt,
				"updated_at": user.UpdatedAt,
			},
		)
	}

	// Return users
	c.JSON(http.StatusOK, response)
}

// UpdateUser handles updating a user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Parse user ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse request body
	var req struct {
		Email string `json:"email" binding:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user
	user, err := h.UserStore.GetByID(id)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get user"},
		)
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update user
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := h.UserStore.Update(user); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to update user"},
		)
		return
	}

	// Return updated user
	c.JSON(
		http.StatusOK, gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	)
}

// DeleteUser handles deleting a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Parse user ID from URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Delete user
	if err := h.UserStore.Delete(id); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to delete user"},
		)
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

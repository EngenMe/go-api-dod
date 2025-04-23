package handlers

import (
	"net/http"

	"github.com/EngenMe/go-api-dod/internal/data/models"
	"github.com/EngenMe/go-api-dod/internal/data/store"
	"github.com/EngenMe/go-api-dod/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler provides handlers for authentication
type AuthHandler struct {
	UserStore      *store.UserStore
	PasswordHasher *utils.PasswordHasher
	TokenManager   *utils.TokenManager
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	userStore *store.UserStore,
	passwordHasher *utils.PasswordHasher,
	tokenManager *utils.TokenManager,
) *AuthHandler {
	return &AuthHandler{
		UserStore:      userStore,
		PasswordHasher: passwordHasher,
		TokenManager:   tokenManager,
	}
}

// Signup handles user registration
func (h *AuthHandler) Signup(c *gin.Context) {
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

	// Hash password
	hashedPassword, err := h.PasswordHasher.Hash(req.Password)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to hash password"},
		)
		return
	}

	// Create user
	user := models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := h.UserStore.Create(&user); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create user"},
		)
		return
	}

	// Generate token
	token, err := h.TokenManager.Generate(user.ID, user.Email)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate token"},
		)
		return
	}

	// Return token
	c.JSON(
		http.StatusCreated, gin.H{
			"token": token,
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
			},
		},
	)
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	// Parse request body
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user
	user, err := h.UserStore.GetByEmail(req.Email)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get user"},
		)
		return
	}
	if user == nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Invalid email or password"},
		)
		return
	}

	// Check password
	if !h.PasswordHasher.Check(req.Password, user.Password) {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Invalid email or password"},
		)
		return
	}

	// Generate token
	token, err := h.TokenManager.Generate(user.ID, user.Email)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate token"},
		)
		return
	}

	// Return token
	c.JSON(
		http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
			},
		},
	)
}

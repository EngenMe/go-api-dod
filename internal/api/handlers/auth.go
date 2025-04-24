package handlers

import (
	"net/http"
	"time"

	"github.com/EngenMe/go-api-dod/internal/data/models"
	"github.com/EngenMe/go-api-dod/internal/data/store"
	"github.com/EngenMe/go-api-dod/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler provides handlers for authentication
type AuthHandler struct {
	UserStore         *store.UserStore
	RefreshTokenStore *store.RefreshTokenStore
	PasswordHasher    *utils.PasswordHasher
	TokenManager      *utils.TokenManager
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	userStore *store.UserStore,
	refreshTokenStore *store.RefreshTokenStore,
	passwordHasher *utils.PasswordHasher,
	tokenManager *utils.TokenManager,
) *AuthHandler {
	return &AuthHandler{
		UserStore:         userStore,
		RefreshTokenStore: refreshTokenStore,
		PasswordHasher:    passwordHasher,
		TokenManager:      tokenManager,
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

	// Check if a user already exists
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

	// Generate an access token
	accessToken, err := h.TokenManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate access token"},
		)
		return
	}

	// Generate refresh token
	refreshTokenString, err := h.TokenManager.GenerateRefreshToken(
		user.ID,
		user.Email,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate refresh token"},
		)
		return
	}

	// Store refresh token in database
	expiresAt := time.Now().Add(h.TokenManager.RefreshTokenExpiresIn)
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: expiresAt,
	}
	if err := h.RefreshTokenStore.Create(refreshToken); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to store refresh token"},
		)
		return
	}

	// Return tokens
	c.JSON(
		http.StatusCreated, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshTokenString,
			"token_type":    "Bearer",
			"expires_in":    int(h.TokenManager.AccessTokenExpiresIn.Seconds()),
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

	// Generate an access token
	accessToken, err := h.TokenManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate access token"},
		)
		return
	}

	// Generate refresh token
	refreshTokenString, err := h.TokenManager.GenerateRefreshToken(
		user.ID,
		user.Email,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate refresh token"},
		)
		return
	}

	// Store refresh token in database
	expiresAt := time.Now().Add(h.TokenManager.RefreshTokenExpiresIn)
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: expiresAt,
	}
	if err := h.RefreshTokenStore.Create(refreshToken); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to store refresh token"},
		)
		return
	}

	// Return tokens
	c.JSON(
		http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshTokenString,
			"token_type":    "Bearer",
			"expires_in":    int(h.TokenManager.AccessTokenExpiresIn.Seconds()),
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
			},
		},
	)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Parse request body
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	claims, err := h.TokenManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Invalid refresh token"},
		)
		return
	}

	// Get refresh token from database
	storedToken, err := h.RefreshTokenStore.GetByToken(req.RefreshToken)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get refresh token"},
		)
		return
	}

	if storedToken == nil || !storedToken.IsValid() {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Invalid or expired refresh token"},
		)
		return
	}

	// Revoke the used refresh token
	if err := h.RefreshTokenStore.Revoke(storedToken.ID); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to revoke refresh token"},
		)
		return
	}

	// Generate a new access token
	accessToken, err := h.TokenManager.GenerateAccessToken(
		claims.UserID,
		claims.Email,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate access token"},
		)
		return
	}

	// Generate a new refresh token
	refreshTokenString, err := h.TokenManager.GenerateRefreshToken(
		claims.UserID,
		claims.Email,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to generate refresh token"},
		)
		return
	}

	// Store a new refresh token in a database
	expiresAt := time.Now().Add(h.TokenManager.RefreshTokenExpiresIn)
	refreshToken := &models.RefreshToken{
		UserID:    claims.UserID,
		Token:     refreshTokenString,
		ExpiresAt: expiresAt,
	}
	if err := h.RefreshTokenStore.Create(refreshToken); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to store refresh token"},
		)
		return
	}

	// Return new tokens
	c.JSON(
		http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshTokenString,
			"token_type":    "Bearer",
			"expires_in":    int(h.TokenManager.AccessTokenExpiresIn.Seconds()),
		},
	)
}

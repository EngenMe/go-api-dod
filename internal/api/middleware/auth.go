package middleware

import (
	"net/http"
	"strings"

	"github.com/EngenMe/go-api-dod/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides authentication middleware for the API
type AuthMiddleware struct {
	TokenManager *utils.TokenManager
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(tokenManager *utils.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		TokenManager: tokenManager,
	}
}

// RequireAuth is a middleware that requires authentication
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Authorization header required"},
			)
			c.Abort()
			return
		}

		// Extract the token from the Authorization header
		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid authorization format"},
			)
			c.Abort()
			return
		}

		// Validate the access token
		claims, err := m.TokenManager.ValidateAccessToken(parts[1])
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid or expired access token"},
			)
			c.Abort()
			return
		}

		// Set user information in the context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

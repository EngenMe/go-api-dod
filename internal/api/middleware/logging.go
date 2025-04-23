package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware provides logging middleware for the API
type LoggingMiddleware struct {
	Logger *log.Logger
}

// NewLoggingMiddleware creates a new LoggingMiddleware
func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{
		Logger: log.New(os.Stdout, "[API] ", log.LstdFlags),
	}
}

// RequestLogger Logger is a middleware that logs API requests
func (m *LoggingMiddleware) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()

		m.Logger.Printf(
			"| %3d | %13v | %15s | %s | %s",
			statusCode, latency, clientIP, method, path,
		)
	}
}

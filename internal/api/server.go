package api

import (
	"github.com/EngenMe/go-api-dod/config"
	"github.com/EngenMe/go-api-dod/internal/api/handlers"
	"github.com/EngenMe/go-api-dod/internal/api/middleware"
	"github.com/EngenMe/go-api-dod/internal/data/store"
	"github.com/EngenMe/go-api-dod/internal/utils"

	"github.com/gin-gonic/gin"
)

// Server represents the API server
type Server struct {
	Router            *gin.Engine
	Config            config.Config
	DB                *store.PostgresStore
	UserStore         *store.UserStore
	RefreshTokenStore *store.RefreshTokenStore
	PasswordHasher    *utils.PasswordHasher
	TokenManager      *utils.TokenManager
	AuthMiddleware    *middleware.AuthMiddleware
	LoggingMiddleware *middleware.LoggingMiddleware
	UserHandler       *handlers.UserHandler
	AuthHandler       *handlers.AuthHandler
}

// NewServer creates a new Server
func NewServer(cfg config.Config, db *store.PostgresStore) *Server {
	// Set Gin mode
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize dependencies
	userStore := store.NewUserStore(db.DB)
	refreshTokenStore := store.NewRefreshTokenStore(db.DB)
	passwordHasher := utils.NewPasswordHasher(cfg.Auth.BcryptCost)
	tokenManager := utils.NewTokenManager(
		cfg.Auth.JWTSecret,
		cfg.Auth.TokenIssuer,
		cfg.Auth.AccessTokenExpiration,
		cfg.Auth.RefreshTokenExpiration,
	)
	authMiddleware := middleware.NewAuthMiddleware(tokenManager)
	loggingMiddleware := middleware.NewLoggingMiddleware()
	userHandler := handlers.NewUserHandler(userStore)
	authHandler := handlers.NewAuthHandler(
		userStore,
		refreshTokenStore,
		passwordHasher,
		tokenManager,
	)

	server := &Server{
		Router:            router,
		Config:            cfg,
		DB:                db,
		UserStore:         userStore,
		RefreshTokenStore: refreshTokenStore,
		PasswordHasher:    passwordHasher,
		TokenManager:      tokenManager,
		AuthMiddleware:    authMiddleware,
		LoggingMiddleware: loggingMiddleware,
		UserHandler:       userHandler,
		AuthHandler:       authHandler,
	}

	// Set up routes
	server.setupRoutes()

	return server
}

// setupRoutes sets up the API routes
func (s *Server) setupRoutes() {
	// Apply middleware
	s.Router.Use(s.LoggingMiddleware.RequestLogger())

	// Versioned API group: /api/v1
	v1 := s.Router.Group("/api/v1")
	{
		// Public routes
		v1.POST("/signup", s.AuthHandler.Signup)
		v1.POST("/login", s.AuthHandler.Login)
		v1.POST("/refresh", s.AuthHandler.RefreshToken)

		// Protected routes
		authorized := v1.Group("/")
		authorized.Use(s.AuthMiddleware.RequireAuth())
		{
			authorized.GET("/users", s.UserHandler.ListUsers)
			authorized.GET("/users/:id", s.UserHandler.GetUser)
			authorized.POST("/users", s.UserHandler.CreateUser)
			authorized.PUT("/users/:id", s.UserHandler.UpdateUser)
			authorized.DELETE("/users/:id", s.UserHandler.DeleteUser)
		}
	}
}

// Run starts the server
func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}

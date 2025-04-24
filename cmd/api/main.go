package main

import (
	"fmt"
	"log"

	"github.com/EngenMe/go-api-dod/config"
	"github.com/EngenMe/go-api-dod/internal/api"
	"github.com/EngenMe/go-api-dod/internal/data/store"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := store.NewPostgresStore(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations (in a production env, use a proper migration tool)
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize and start an API server
	server := api.NewServer(cfg, db)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := server.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

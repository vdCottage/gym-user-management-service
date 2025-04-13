package main

import (
	"log"

	"github.com/yourname/fitness-platform/internal/config"
	"github.com/yourname/fitness-platform/internal/database"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Run migrations
	if err := database.Migrate(cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}

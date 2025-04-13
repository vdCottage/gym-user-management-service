package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourname/fitness-platform/pkg/logger"
	"gorm.io/gorm"
)

// RunMigrations runs all SQL migration files in the migrations directory
func RunMigrations(db *gorm.DB, log *logger.Logger) error {
	migrationsDir := "internal/database/migrations"

	// Read all files in the migrations directory
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort files by name to ensure they run in order
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			// Read the migration file
			content, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
			}

			// Execute the migration
			log.Info("Running migration", map[string]interface{}{
				"file": file.Name(),
			})

			if err := db.Exec(string(content)).Error; err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
			}

			log.Info("Migration completed", map[string]interface{}{
				"file": file.Name(),
			})
		}
	}

	return nil
}

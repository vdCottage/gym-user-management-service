package database

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/yourname/fitness-platform/internal/config"
	"github.com/yourname/fitness-platform/internal/models"
	appLogger "github.com/yourname/fitness-platform/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establishes a connection to the database
func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	// Configure GORM logger
	gormLogger := gormLogger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		gormLogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormLogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	// Ensure uuid extension is enabled
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.GymOwner{},
		&models.Trainer{},
		&models.Customer{},
		&models.OTP{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	log := appLogger.NewLogger("info", "json", "")
	log.Info("Database connection established", map[string]interface{}{
		"host": cfg.Database.Host,
		"port": cfg.Database.Port,
		"name": cfg.Database.DBName,
	})

	return nil
}

// Migrate runs database migrations
func Migrate(cfg *config.Config) error {
	if DB == nil {
		if err := Connect(cfg); err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	driver, err := pgmigrate.WithInstance(sqlDB, &pgmigrate.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	return sqlDB.Close()
}

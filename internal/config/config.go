package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Auth        AuthConfig
	Redis       RedisConfig
	Log         LogConfig
	JWT         JWTConfig
	Env         string
	OTP         struct {
		Length        int    `mapstructure:"length"`
		ExpiryMinutes int    `mapstructure:"expiry_minutes"`
		RateLimit     int    `mapstructure:"rate_limit"`
		DefaultOTP    string `mapstructure:"default_otp"`
	} `mapstructure:"otp"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret           string
	JWTExpirationHours  int
	OTPExpirationMins   int
	PasswordMinLength   int
	MaxLoginAttempts    int
	LockoutDurationMins int
	DefaultOTP          string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level      string
	Format     string
	OutputPath string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{
		Environment: getEnv("APP_ENV", "development"),
		Env:         getEnv("APP_ENV", "development"),
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "fitness_platform"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret:           getEnv("JWT_SECRET", "your-secret-key"),
			JWTExpirationHours:  getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
			OTPExpirationMins:   getEnvAsInt("OTP_EXPIRATION_MINS", 5),
			PasswordMinLength:   getEnvAsInt("PASSWORD_MIN_LENGTH", 8),
			MaxLoginAttempts:    getEnvAsInt("MAX_LOGIN_ATTEMPTS", 5),
			LockoutDurationMins: getEnvAsInt("LOCKOUT_DURATION_MINS", 15),
			DefaultOTP:          getEnv("DEFAULT_OTP", "123456"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			OutputPath: getEnv("LOG_OUTPUT_PATH", ""),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		},
	}

	// Set default OTP values
	config.OTP.Length = getEnvAsIntOrDefault("OTP_LENGTH", 6)
	config.OTP.ExpiryMinutes = getEnvAsIntOrDefault("OTP_EXPIRY_MINUTES", 5)
	config.OTP.RateLimit = getEnvAsIntOrDefault("OTP_RATE_LIMIT", 3)
	config.OTP.DefaultOTP = getEnvOrDefault("OTP_DEFAULT", "123456")

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	return nil
}

// Helper functions to get environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

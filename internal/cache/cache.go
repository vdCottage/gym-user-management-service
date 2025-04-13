package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yourname/fitness-platform/config"
	"github.com/yourname/fitness-platform/pkg/logger"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

// Connect establishes a connection to Redis
func Connect(cfg *config.Config) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test the connection
	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log := logger.NewLogger("info", "json", "")
	log.Info("Redis connection established", map[string]interface{}{
		"host": cfg.Redis.Host,
		"port": cfg.Redis.Port,
		"db":   cfg.Redis.DB,
	})

	return nil
}

// Close closes the Redis connection
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// Set stores a key-value pair in Redis
func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(Ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis by key
func Get(key string) (string, error) {
	return Client.Get(Ctx, key).Result()
}

// Delete removes a key from Redis
func Delete(key string) error {
	return Client.Del(Ctx, key).Err()
}

// Exists checks if a key exists in Redis
func Exists(key string) (bool, error) {
	result, err := Client.Exists(Ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

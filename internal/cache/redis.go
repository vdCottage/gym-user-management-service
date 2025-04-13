package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yourname/fitness-platform/internal/config"
	"github.com/yourname/fitness-platform/pkg/logger"
)

// RedisClient defines the interface for Redis operations
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

// redisClient implements the RedisClient interface
type redisClient struct {
	client *redis.Client
	logger *logger.Logger
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.RedisConfig, log *logger.Logger) (RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Redis connection established", map[string]interface{}{
		"host": cfg.Host,
		"port": cfg.Port,
		"db":   cfg.DB,
	})

	return &redisClient{
		client: client,
		logger: log,
	}, nil
}

// Set stores a key-value pair in Redis
func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

// Get retrieves a value from Redis by key
func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return data, nil
}

// Delete removes a key from Redis
func (r *redisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (r *redisClient) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// Incr increments the value of a key in Redis
func (r *redisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Expire sets the expiration time for a key in Redis
func (r *redisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// SetOTP stores an OTP in Redis
func (r *redisClient) SetOTP(ctx context.Context, key string, code string, expiry time.Duration) error {
	return r.client.Set(ctx, key, code, expiry).Err()
}

// GetOTP retrieves an OTP from Redis
func (r *redisClient) GetOTP(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

// DeleteOTP deletes an OTP from Redis
func (r *redisClient) DeleteOTP(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

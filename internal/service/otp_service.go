package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/yourname/fitness-platform/internal/cache"
	"github.com/yourname/fitness-platform/internal/config"
	"github.com/yourname/fitness-platform/internal/models"
	"github.com/yourname/fitness-platform/internal/repository"
)

var (
	ErrInvalidOTP        = fmt.Errorf("invalid OTP")
	ErrOTPExpired        = fmt.Errorf("OTP has expired")
	ErrRateLimitExceeded = fmt.Errorf("rate limit exceeded")
)

// RateLimiter interface for rate limiting
type RateLimiter interface {
	IsAllowed(key string) bool
	Increment(key string)
}

// redisRateLimiter implements RateLimiter using Redis
type redisRateLimiter struct {
	client cache.RedisClient
	prefix string
	limit  int
	window time.Duration
}

// NewRedisRateLimiter creates a new Redis-based rate limiter
func NewRedisRateLimiter(client cache.RedisClient, prefix string, limit int, window time.Duration) RateLimiter {
	return &redisRateLimiter{
		client: client,
		prefix: prefix,
		limit:  limit,
		window: window,
	}
}

// IsAllowed checks if a request is allowed under rate limiting
func (r *redisRateLimiter) IsAllowed(key string) bool {
	count, err := r.client.Get(context.Background(), r.prefix+key)
	if err != nil {
		return true // Allow if Redis is down
	}

	countInt, _ := strconv.Atoi(count)
	return countInt < r.limit
}

// Increment increments the rate limit counter
func (r *redisRateLimiter) Increment(key string) {
	ctx := context.Background()
	_, _ = r.client.Incr(ctx, r.prefix+key)
	_ = r.client.Expire(ctx, r.prefix+key, r.window)
}

// OTPService handles OTP-related business logic
type OTPService struct {
	config      *config.Config
	redisClient cache.RedisClient
	userRepo    *repository.CustomerRepository
	trainerRepo *repository.TrainerRepository
	logger      *zap.Logger
}

// NewOTPService creates a new OTP service
func NewOTPService(cfg *config.Config, redisClient cache.RedisClient, userRepo *repository.CustomerRepository, trainerRepo *repository.TrainerRepository, logger *zap.Logger) *OTPService {
	return &OTPService{
		config:      cfg,
		redisClient: redisClient,
		userRepo:    userRepo,
		trainerRepo: trainerRepo,
		logger:      logger,
	}
}

// GenerateOTP generates a new OTP for the given target
func (s *OTPService) GenerateOTP(ctx context.Context, target string, targetType models.OTPTarget) (string, error) {
	// Check rate limit
	rateLimitKey := fmt.Sprintf("otp:ratelimit:%s", target)
	exists, err := s.redisClient.Exists(ctx, rateLimitKey)
	if err != nil {
		s.logger.Error("failed to check rate limit",
			zap.String("target", target),
			zap.Error(err))
		return "", err
	}
	if exists {
		return "", models.ErrOTPExpired
	}

	// Generate OTP
	otp := s.generateRandomOTP()

	// Store OTP in Redis
	otpKey := fmt.Sprintf("otp:%s:%s", targetType, target)
	err = s.redisClient.Set(ctx, otpKey, otp, time.Duration(s.config.OTP.ExpiryMinutes)*time.Minute)
	if err != nil {
		s.logger.Error("failed to store OTP",
			zap.String("target", target),
			zap.Error(err))
		return "", err
	}

	// Set rate limit
	err = s.redisClient.Set(ctx, rateLimitKey, "1", time.Duration(s.config.OTP.RateLimit)*time.Minute)
	if err != nil {
		s.logger.Error("failed to set rate limit",
			zap.String("target", target),
			zap.Error(err))
		return "", err
	}

	return otp, nil
}

// VerifyOTP verifies the OTP for the given target
func (s *OTPService) VerifyOTP(ctx context.Context, target string, targetType models.OTPTarget, otp string) (bool, error) {
	otpKey := fmt.Sprintf("otp:%s:%s", targetType, target)
	storedOTP, err := s.redisClient.Get(ctx, otpKey)
	if err != nil {
		s.logger.Error("failed to get OTP",
			zap.String("target", target),
			zap.Error(err))
		return false, err
	}

	if storedOTP == "" {
		return false, models.ErrOTPExpired
	}

	if storedOTP != otp {
		return false, models.ErrInvalidOTP
	}

	// Delete OTP after successful verification
	err = s.redisClient.Delete(ctx, otpKey)
	if err != nil {
		s.logger.Error("failed to delete OTP",
			zap.String("target", target),
			zap.Error(err))
		return false, err
	}

	// Activate user based on target type
	userID, err := uuid.Parse(target)
	if err != nil {
		s.logger.Error("failed to parse user ID",
			zap.String("target", target),
			zap.Error(err))
		return false, err
	}

	if err := s.activateUser(ctx, userID, targetType); err != nil {
		s.logger.Error("failed to activate user",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return false, err
	}

	return true, nil
}

// generateRandomOTP generates a random OTP of the configured length
func (s *OTPService) generateRandomOTP() string {
	if s.config.OTP.DefaultOTP != "" {
		return s.config.OTP.DefaultOTP
	}

	digits := make([]byte, s.config.OTP.Length)
	for i := range digits {
		digits[i] = byte(rand.Intn(10) + '0')
	}
	return string(digits)
}

// activateUser marks a user as active based on their type
func (s *OTPService) activateUser(ctx context.Context, userID uuid.UUID, targetType models.OTPTarget) error {
	switch targetType {
	case models.OTPTargetTrainer:
		trainer, err := s.trainerRepo.GetByID(ctx, userID)
		if err != nil {
			return err
		}
		trainer.IsActive = true
		return s.trainerRepo.Update(ctx, trainer)

	default:
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return err
		}
		user.IsActive = true
		return s.userRepo.Update(ctx, user)
	}
}

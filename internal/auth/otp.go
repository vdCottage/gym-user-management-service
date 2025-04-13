package auth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/yourname/fitness-platform/config"
	"github.com/yourname/fitness-platform/internal/cache"
	"github.com/yourname/fitness-platform/internal/database"
	"github.com/yourname/fitness-platform/internal/models"
	"github.com/yourname/fitness-platform/pkg/logger"
)

// GenerateOTP generates a new OTP for a user
func GenerateOTP(userID uuid.UUID, otpType models.OTPType, cfg *config.Config) (string, error) {
	// Generate a random 6-digit code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Create OTP record in database
	otp := &models.OTP{
		UserID:    userID,
		Code:      code,
		Type:      otpType,
		ExpiresAt: time.Now().Add(cfg.OTP.Expiration),
	}

	// Save OTP to database
	if err := database.DB.Create(otp).Error; err != nil {
		return "", fmt.Errorf("failed to save OTP: %w", err)
	}

	// Store OTP in Redis for faster verification
	key := fmt.Sprintf("otp:%s:%s", userID.String(), otpType)
	if err := cache.Set(key, code, cfg.OTP.Expiration); err != nil {
		log := logger.NewLogger("warn", "json", "")
		log.Warn("Failed to store OTP in cache", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return code, nil
}

// VerifyOTP verifies an OTP for a user
func VerifyOTP(userID uuid.UUID, code string, otpType models.OTPType) (bool, error) {
	// Check Redis first for faster verification
	key := fmt.Sprintf("otp:%s:%s", userID.String(), otpType)
	cachedCode, err := cache.Get(key)
	if err == nil && cachedCode == code {
		// OTP verified from cache, mark as used in database
		if err := markOTPAsUsed(userID, code, otpType); err != nil {
			log := logger.NewLogger("warn", "json", "")
			log.Warn("Failed to mark OTP as used", map[string]interface{}{
				"error": err.Error(),
			})
		}
		return true, nil
	}

	// Check database if not found in cache
	var otp models.OTP
	if err := database.DB.Where("user_id = ? AND code = ? AND type = ? AND used = ? AND expires_at > ?",
		userID, code, otpType, false, time.Now()).First(&otp).Error; err != nil {
		return false, fmt.Errorf("invalid or expired OTP")
	}

	// Mark OTP as used
	if err := markOTPAsUsed(userID, code, otpType); err != nil {
		log := logger.NewLogger("warn", "json", "")
		log.Warn("Failed to mark OTP as used", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return true, nil
}

// markOTPAsUsed marks an OTP as used in the database
func markOTPAsUsed(userID uuid.UUID, code string, otpType models.OTPType) error {
	// Update OTP in database
	if err := database.DB.Model(&models.OTP{}).
		Where("user_id = ? AND code = ? AND type = ?", userID, code, otpType).
		Update("used", true).Error; err != nil {
		return fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	// Delete OTP from Redis
	key := fmt.Sprintf("otp:%s:%s", userID.String(), otpType)
	if err := cache.Delete(key); err != nil {
		log := logger.NewLogger("warn", "json", "")
		log.Warn("Failed to delete OTP from cache", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return nil
}

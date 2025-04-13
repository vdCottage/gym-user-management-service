package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourname/fitness-platform/internal/models"
	"gorm.io/gorm"
)

var (
	ErrOTPNotFound = errors.New("OTP not found")
	ErrOTPExpired  = errors.New("OTP has expired")
	ErrOTPUsed     = errors.New("OTP has already been used")
)

// OTPRepository handles OTP data persistence
type OTPRepository struct {
	db *gorm.DB
}

// NewOTPRepository creates a new OTPRepository
func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{
		db: db,
	}
}

// Create creates a new OTP
func (r *OTPRepository) Create(ctx context.Context, otp *models.OTP) error {
	return r.db.WithContext(ctx).Create(otp).Error
}

// GetByID retrieves an OTP by ID
func (r *OTPRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.OTP, error) {
	var otp models.OTP
	err := r.db.WithContext(ctx).First(&otp, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOTPNotFound
		}
		return nil, err
	}
	return &otp, nil
}

// GetByUserIDAndType retrieves an OTP by user ID and type
func (r *OTPRepository) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, otpType models.OTPType) (*models.OTP, error) {
	var otp models.OTP
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND used = ? AND expires_at > ?",
			userID, otpType, false, time.Now()).
		Order("created_at DESC").
		First(&otp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOTPNotFound
		}
		return nil, err
	}
	return &otp, nil
}

// MarkAsUsed marks an OTP as used
func (r *OTPRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.OTP{}).Where("id = ?", id).Update("used", true).Error
}

// DeleteExpired deletes expired OTPs
func (r *OTPRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ? OR used = ?", time.Now(), true).Delete(&models.OTP{}).Error
}

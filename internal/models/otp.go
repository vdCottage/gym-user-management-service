package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidOTP = errors.New("invalid OTP")
	ErrOTPExpired = errors.New("OTP has expired")
)

// OTPType represents the type of OTP
type OTPType string

const (
	OTPTypeEmail OTPType = "email"
	OTPTypePhone OTPType = "phone"
)

// OTPTarget represents the target user type for OTP
type OTPTarget string

const (
	OTPTargetGymOwner OTPTarget = "gym_owner"
	OTPTargetTrainer  OTPTarget = "trainer"
	OTPTargetCustomer OTPTarget = "customer"
)

// OTP represents a one-time password
// @Description One-time password information
type OTP struct {
	// @Description Unique identifier for the OTP
	// @example 123e4567-e89b-12d3-a456-426614174000
	ID uuid.UUID `json:"-" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" validate:"required,uuid4"`

	// @Description Target user ID
	// @example 123e4567-e89b-12d3-a456-426614174000
	UserID uuid.UUID `json:"-" gorm:"type:uuid;not null;index" validate:"required,uuid4"`

	// @Description Target type (email or phone)
	// @example email
	Type OTPType `json:"type" gorm:"type:string;not null" validate:"required,oneof=email phone"`

	// @Description Target value (email or phone number)
	// @example john.doe@example.com
	Target string `json:"target" gorm:"not null" validate:"required"`

	// @Description Target user type
	// @example customer
	TargetType OTPTarget `json:"target_type" gorm:"type:string;not null" validate:"required,oneof=gym_owner trainer customer"`

	// @Description OTP code
	// @example 123456
	Code string `json:"-" gorm:"not null" validate:"required,len=6"`

	// @Description Whether the OTP has been used
	// @example false
	Used bool `json:"used" gorm:"not null;default:false"`

	// @Description When the OTP expires
	// @example 2024-04-12T12:00:00Z
	ExpiresAt time.Time `json:"expires_at" gorm:"not null" validate:"required"`

	// @Description Creation timestamp
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// @Description Last update timestamp
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for the OTP model
func (OTP) TableName() string {
	return "otps"
}

// IsExpired checks if the OTP has expired
func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// IsValid checks if the OTP is valid (not expired and not used)
func (o *OTP) IsValid() bool {
	return !o.IsExpired() && !o.Used
}

// SendOTPRequest represents a request to send OTP
// @Description Request to send OTP
type SendOTPRequest struct {
	// @Description Target type (email or phone)
	// @example email
	Type OTPType `json:"type" validate:"required,oneof=email phone"`

	// @Description Target value (email or phone number)
	// @example john.doe@example.com
	Target string `json:"target" validate:"required"`

	// @Description Target user type
	// @example customer
	TargetType OTPTarget `json:"target_type" validate:"required,oneof=gym_owner trainer customer"`
}

// Note: VerifyOTPRequest is now defined in internal/models/auth.go

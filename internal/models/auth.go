package models

// UserType represents the type of user
type UserType string

const (
	UserTypeGymOwner UserType = "gym_owner"
	UserTypeTrainer  UserType = "trainer"
	UserTypeCustomer UserType = "customer"
)

// LoginRequest represents a login request
// @Description Login request details
type LoginRequest struct {
	// @Description Email address or phone number
	// @example john.doe@example.com
	Username string `json:"username" validate:"required"`

	// @Description Password
	// @example securePassword123
	Password string `json:"password" validate:"required,min=8"`

	// @Description Type of user (gym_owner, trainer, customer)
	// @example trainer
	UserType UserType `json:"user_type" validate:"required,oneof=gym_owner trainer customer"`
}

// LoginResponse represents a login response
// @Description Login response details
type LoginResponse struct {
	// @Description Access token
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	AccessToken string `json:"access_token"`

	// @Description Token type
	// @example Bearer
	TokenType string `json:"token_type"`

	// @Description User ID
	// @example 123e4567-e89b-12d3-a456-426614174000
	UserID string `json:"user_id"`

	// @Description User type
	// @example trainer
	UserType UserType `json:"user_type"`
}

// VerifyOTPRequest represents an OTP verification request
// @Description OTP verification request details
type VerifyOTPRequest struct {
	// @Description User ID
	// @example 123e4567-e89b-12d3-a456-426614174000
	UserID string `json:"user_id" validate:"required,uuid"`

	// @Description OTP code
	// @example 123456
	OTP string `json:"otp" validate:"required,len=6"`

	// @Description Type of user (gym_owner, trainer, customer)
	// @example trainer
	UserType UserType `json:"user_type" validate:"required,oneof=gym_owner trainer customer"`
}

// VerifyOTPResponse represents an OTP verification response
// @Description OTP verification response details
type VerifyOTPResponse struct {
	// @Description Access token
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	AccessToken string `json:"access_token"`

	// @Description Token type
	// @example Bearer
	TokenType string `json:"token_type"`

	// @Description User ID
	// @example 123e4567-e89b-12d3-a456-426614174000
	UserID string `json:"user_id"`

	// @Description User type
	// @example trainer
	UserType UserType `json:"user_type"`

	// @Description Whether the user is active
	// @example true
	IsActive bool `json:"is_active"`
}

// RefreshTokenRequest represents a token refresh request
// @Description Token refresh request details
type RefreshTokenRequest struct {
	// @Description Refresh token
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents a token refresh response
// @Description Token refresh response details
type RefreshTokenResponse struct {
	// @Description Access token
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	AccessToken string `json:"access_token"`

	// @Description Token type
	// @example Bearer
	TokenType string `json:"token_type"`
}

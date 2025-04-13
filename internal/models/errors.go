package models

import "errors"

// Common errors
var (
	ErrNotFound       = errors.New("not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInternalServer = errors.New("internal server error")
)

// OTP errors
var (
	// ErrOTPExpired           = errors.New("otp expired")
	ErrOTPInvalid           = errors.New("invalid otp")
	ErrOTPRateLimitExceeded = errors.New("otp rate limit exceeded")
	ErrOTPNotFound          = errors.New("otp not found")
)

// User errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

package handlers

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yourname/fitness-platform/internal/models"
	"github.com/yourname/fitness-platform/internal/service"
	"go.uber.org/zap"
)

// OTPHandler handles OTP-related HTTP requests
type OTPHandler struct {
	otpService *service.OTPService
	logger     *zap.Logger
	validator  *validator.Validate
}

// NewOTPHandler creates a new OTPHandler
func NewOTPHandler(otpService *service.OTPService, logger *zap.Logger) *OTPHandler {
	return &OTPHandler{
		otpService: otpService,
		logger:     logger,
		validator:  validator.New(),
	}
}

// @Summary Send OTP
// @Description Send an OTP to the specified target (email/phone)
// @Tags otp
// @Accept json
// @Produce json
// @Param request body models.SendOTPRequest true "OTP Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /otp/send [post]
func (h *OTPHandler) SendOTP(c *fiber.Ctx) error {
	var req models.SendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	otp, err := h.otpService.GenerateOTP(c.Context(), req.Target, req.TargetType)
	if err != nil {
		h.logger.Error("Failed to generate OTP",
			zap.Error(err),
			zap.String("target", req.Target),
			zap.String("type", string(req.Type)))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate OTP",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OTP sent successfully",
		"otp":     otp, // Only included in development environment
	})
}

// @Summary Verify OTP
// @Description Verify an OTP for a user
// @Tags otp
// @Accept json
// @Produce json
// @Param request body models.VerifyOTPRequest true "Verify OTP Request"
// @Success 200 {object} models.VerifyOTPResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /otp/verify [post]
func (h *OTPHandler) VerifyOTP(c *fiber.Ctx) error {
	var req models.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Parse UUID from string
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Convert UserType to OTPTarget
	targetType := models.OTPTarget(req.UserType)

	verified, err := h.otpService.VerifyOTP(c.Context(), req.UserID, targetType, req.OTP)
	if err != nil {
		h.logger.Error("Failed to verify OTP",
			zap.Error(err),
			zap.String("user_id", userID.String()))

		if errors.Is(err, models.ErrOTPExpired) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "OTP has expired",
			})
		}
		if errors.Is(err, models.ErrOTPInvalid) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid OTP",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify OTP",
		})
	}

	if !verified {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid OTP",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "OTP verified successfully",
		"user_id":   userID.String(),
		"user_type": req.UserType,
		"is_active": true,
	})
}

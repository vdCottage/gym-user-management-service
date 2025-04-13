package response

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Response represents a standardized API response
type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data,omitempty"`
}

// Success sends a successful response
func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success:    true,
		Message:    message,
		StatusCode: fiber.StatusOK,
		Data:       data,
	})
}

// Created sends a 201 Created response
func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success:    true,
		Message:    message,
		StatusCode: fiber.StatusCreated,
		Data:       data,
	})
}

// Error sends an error response
func Error(c *fiber.Ctx, statusCode int, message string, err error, logger *zap.Logger) error {
	// Log the error with appropriate level based on status code
	if statusCode >= 500 {
		logger.Error("Server error",
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.Int("status", statusCode))
	} else {
		logger.Warn("Client error",
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.Int("status", statusCode))
	}

	return c.Status(statusCode).JSON(Response{
		Success:    false,
		Message:    message,
		StatusCode: statusCode,
	})
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *fiber.Ctx, message string, err error, logger *zap.Logger) error {
	return Error(c, fiber.StatusBadRequest, message, err, logger)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *fiber.Ctx, message string, err error, logger *zap.Logger) error {
	return Error(c, fiber.StatusUnauthorized, message, err, logger)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *fiber.Ctx, message string, err error, logger *zap.Logger) error {
	return Error(c, fiber.StatusForbidden, message, err, logger)
}

// NotFound sends a 404 Not Found response
func NotFound(c *fiber.Ctx, message string, err error, logger *zap.Logger) error {
	return Error(c, fiber.StatusNotFound, message, err, logger)
}

// ServerError sends a 500 Internal Server Error response
func ServerError(c *fiber.Ctx, message string, err error, logger *zap.Logger) error {
	return Error(c, fiber.StatusInternalServerError, message, err, logger)
}

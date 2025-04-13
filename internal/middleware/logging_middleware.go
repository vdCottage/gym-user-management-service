package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// LoggingMiddleware creates a middleware that logs HTTP requests
func LoggingMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Get request body for logging
		reqBody := string(c.Body())

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get error message if any
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		// Log request details
		logger.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", duration),
			zap.String("ip", c.IP()),
			zap.String("request_body", reqBody),
			zap.String("error", errMsg),
		)

		return err
	}
}

package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yourname/fitness-platform/pkg/logger"
)

// Logger middleware for Fiber
func Logger(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get status code
		status := c.Response().StatusCode()

		// Log request
		log.Info("HTTP Request", map[string]interface{}{
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     status,
			"duration":   duration.Milliseconds(),
			"ip":         c.IP(),
			"user_agent": c.Get("User-Agent"),
		})

		return err
	}
}

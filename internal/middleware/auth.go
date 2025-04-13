package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yourname/fitness-platform/internal/auth"
	"github.com/yourname/fitness-platform/internal/config"
	"github.com/yourname/fitness-platform/internal/models"
)

// FiberAuthMiddleware creates a new authentication middleware for Fiber
func FiberAuthMiddleware(cfg *config.Config) fiber.Handler {
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpirationHours)

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_role", claims.Role)

		return c.Next()
	}
}

// RequireRole creates a middleware that checks if the user has the required role
func RequireRole(roles ...models.UserType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role").(models.UserType)
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}

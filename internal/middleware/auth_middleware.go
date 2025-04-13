package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yourname/fitness-platform/internal/auth"
	"github.com/yourname/fitness-platform/internal/models"
)

// AuthMiddleware creates a middleware for JWT authentication
func AuthMiddleware(jwtService *auth.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Extract token
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		// Store user claims in context
		c.Locals("userID", claims.UserID)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// RoleAuthMiddleware creates a middleware for role-based authorization
func RoleAuthMiddleware(jwtService *auth.JWTService, allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Extract token
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		// Check if user role is allowed
		userRole := claims.Role
		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "insufficient permissions",
			})
		}

		// Store user claims in context
		c.Locals("userID", claims.UserID)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// GymOwnerAuthMiddleware middleware specifically for gym owner routes
func GymOwnerAuthMiddleware(jwtService *auth.JWTService) fiber.Handler {
	return RoleAuthMiddleware(jwtService, models.RoleGymOwner)
}

// TrainerAuthMiddleware middleware specifically for trainer routes
func TrainerAuthMiddleware(jwtService *auth.JWTService) fiber.Handler {
	return RoleAuthMiddleware(jwtService, models.RoleTrainer)
}

// CustomerAuthMiddleware middleware specifically for customer routes
func CustomerAuthMiddleware(jwtService *auth.JWTService) fiber.Handler {
	return RoleAuthMiddleware(jwtService, models.RoleCustomer)
}

// AdminOrGymOwnerAuthMiddleware middleware for routes accessible by both admin and gym owner
func AdminOrGymOwnerAuthMiddleware(jwtService *auth.JWTService) fiber.Handler {
	return RoleAuthMiddleware(jwtService, models.RoleAdmin, models.RoleGymOwner)
}

package middleware

// import (
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/google/uuid"
// 	"github.com/yourname/fitness-platform/internal/auth"
// 	"go.uber.org/zap"
// )

// // Chain applies middlewares in order
// func Chain(app *fiber.App, middlewares ...fiber.Handler) {
// 	for _, middleware := range middlewares {
// 		app.Use(middleware)
// 	}
// }

// // AuthMiddleware verifies JWT token and adds user info to context
// func AuthhMiddleware(jwtService *auth.JWTService) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		token := c.Get("Authorization")
// 		if token == "" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Unauthorized",
// 			})
// 		}

// 		// Remove "Bearer " prefix if present
// 		if len(token) > 7 && token[:7] == "Bearer " {
// 			token = token[7:]
// 		}

// 		claims, err := jwtService.ValidateToken(token)
// 		if err != nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Invalid token",
// 			})
// 		}

// 		// Add claims to context
// 		c.Locals("user_id", claims.UserID)
// 		c.Locals("user_role", claims.Role)
// 		c.Locals("user_email", claims.Email)

// 		return c.Next()
// 	}
// }

// // LoggingMiddleware logs request details
// func LoggingMiddleware(log *zap.Logger) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		start := time.Now()

// 		// Process request
// 		err := c.Next()

// 		// Log request details
// 		duration := time.Since(start)
// 		log.Info("request completed",
// 			zap.String("method", c.Method()),
// 			zap.String("path", c.Path()),
// 			zap.Int("status", c.Response().StatusCode()),
// 			zap.Duration("duration", duration),
// 			zap.String("ip", c.IP()),
// 			zap.String("user_agent", c.Get("User-Agent")),
// 		)

// 		return err
// 	}
// }

// // CORSMiddleware adds CORS headers
// func CORSMiddleware() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		c.Set("Access-Control-Allow-Origin", "*")
// 		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 		if c.Method() == "OPTIONS" {
// 			return c.SendStatus(fiber.StatusOK)
// 		}

// 		return c.Next()
// 	}
// }

// // RecoveryMiddleware recovers from panics
// func RecoveryMiddleware(log *zap.Logger) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				err, ok := r.(error)
// 				if !ok {
// 					err = fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
// 				}

// 				log.Error("panic recovered",
// 					zap.Any("error", err),
// 					zap.String("path", c.Path()),
// 					zap.String("method", c.Method()),
// 					zap.String("ip", c.IP()),
// 				)

// 				// Check if response has already been written
// 				if c.Response().StatusCode() == 0 {
// 					c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 						"error": "Internal Server Error",
// 					})
// 				}
// 			}
// 		}()

// 		return c.Next()
// 	}
// }

// // RateLimitMiddleware implements rate limiting
// func RateLimitMiddleware(requests int, duration time.Duration) fiber.Handler {
// 	// Create a map to store IP addresses and their request counts
// 	type client struct {
// 		count    int
// 		lastSeen time.Time
// 	}
// 	clients := make(map[string]*client)

// 	// Start a goroutine to clean up old entries
// 	go func() {
// 		for {
// 			time.Sleep(duration)
// 			// Clean up old entries
// 			now := time.Now()
// 			for ip, c := range clients {
// 				if now.Sub(c.lastSeen) > duration {
// 					delete(clients, ip)
// 				}
// 			}
// 		}
// 	}()

// 	return func(c *fiber.Ctx) error {
// 		ip := c.IP()
// 		now := time.Now()

// 		// Get or create client
// 		cl, exists := clients[ip]
// 		if !exists {
// 			cl = &client{count: 0, lastSeen: now}
// 			clients[ip] = cl
// 		}

// 		// Reset count if duration has passed
// 		if now.Sub(cl.lastSeen) > duration {
// 			cl.count = 0
// 			cl.lastSeen = now
// 		}

// 		// Check if rate limit exceeded
// 		if cl.count >= requests {
// 			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
// 				"error": "Rate limit exceeded",
// 			})
// 		}

// 		// Increment count and update last seen
// 		cl.count++
// 		cl.lastSeen = now

// 		return c.Next()
// 	}
// }

// // RequestIDMiddleware adds a unique request ID to each request
// func RequestIDMiddleware() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		requestID := c.Get("X-Request-ID")
// 		if requestID == "" {
// 			// Generate a new request ID using UUID
// 			requestID = uuid.New().String()
// 		}
// 		c.Locals("request_id", requestID)
// 		c.Set("X-Request-ID", requestID)
// 		return c.Next()
// 	}
// }

// // TimeoutMiddleware adds a timeout to the request context
// func TimeoutMiddleware(timeout time.Duration) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// Create a channel to signal completion
// 		done := make(chan bool)

// 		// Start a goroutine to process the request
// 		go func() {
// 			_ = c.Next()
// 			done <- true
// 		}()

// 		// Wait for either completion or timeout
// 		select {
// 		case <-done:
// 			return nil
// 		case <-time.After(timeout):
// 			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
// 				"error": "Request timeout",
// 			})
// 		}
// 	}
// }

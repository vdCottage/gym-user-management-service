package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/yourname/fitness-platform/docs"
	"github.com/yourname/fitness-platform/internal/auth"
	"github.com/yourname/fitness-platform/internal/cache"
	"github.com/yourname/fitness-platform/internal/config"
	"github.com/yourname/fitness-platform/internal/database"
	"github.com/yourname/fitness-platform/internal/handlers"
	"github.com/yourname/fitness-platform/internal/middleware"
	"github.com/yourname/fitness-platform/internal/repository"
	"github.com/yourname/fitness-platform/internal/service"
	"github.com/yourname/fitness-platform/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title Fitness Platform API
// @version 2.0
// @description A comprehensive fitness platform API that manages gym owners, trainers, and customers with features including user management, authentication, profile management, and role-based access control.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.fitnessplatform.io/support
// @contact.email support@fitnessplatform.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token for authentication. Use the format: Bearer <token>

// @tag.name gym-owners
// @tag.description Operations about gym owners including registration, authentication, profile management, and trainer/customer management. Gym owners can manage their trainers and customers.

// @tag.name trainers
// @tag.description Operations about trainers including registration, profile management, and availability. Trainers are managed by gym owners and can view customer information.

// @tag.name customers
// @tag.description Operations about customers including registration, profile management, fitness goals, and health tracking. Customers can be managed by gym owners and viewed by trainers.

// @tag.name otp
// @tag.description Operations for One-Time Password (OTP) verification. Used for account verification and password reset.

// @tag.name health
// @tag.description Health check endpoint for monitoring service status and API availability.

// @x-request-id {"type": "string", "in": "header", "name": "Request-Id", "description": "Unique identifier for tracking requests across the system"}

// @x-success {"codes": [200, 201, 204], "descriptions": {"200": "Success", "201": "Created", "204": "No Content"}}
// @x-errors {"codes": [400, 401, 403, 404, 409, 422, 500], "descriptions": {"400": "Bad Request", "401": "Unauthorized", "403": "Forbidden", "404": "Not Found", "409": "Conflict", "422": "Unprocessable Entity", "500": "Internal Server Error"}}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.TimeKey = "timestamp"
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapLogger, err := logConfig.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	appLogger := logger.NewLogger("info", "json", "")

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Initialize Redis client
	redisClient, err := cache.NewRedisClient(&cfg.Redis, appLogger)
	if err != nil {
		zapLogger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 10,
	})

	// Add middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.LoggingMiddleware(zapLogger))

	// Initialize JWT service
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpirationHours)

	// Initialize repositories
	trainerRepo := repository.NewTrainerRepository(database.DB)
	gymOwnerRepo := repository.NewGymOwnerRepository(database.DB)
	customerRepo := repository.NewCustomerRepository(database.DB)

	// Initialize services
	otpService := service.NewOTPService(cfg, redisClient, customerRepo, trainerRepo, zapLogger)
	trainerService := service.NewTrainerService(trainerRepo, cfg, appLogger, otpService)
	gymOwnerService, err := service.NewGymOwnerService(
		gymOwnerRepo,
		trainerRepo,
		customerRepo,
		redisClient,
		cfg,
		zapLogger,
	)
	if err != nil {
		zapLogger.Fatal("Failed to create gym owner service", zap.Error(err))
	}
	customerService, err := service.NewCustomerService(customerRepo, redisClient, cfg, zapLogger)
	if err != nil {
		zapLogger.Fatal("Failed to create customer service", zap.Error(err))
	}

	// Initialize handlers
	trainerHandler := handlers.NewTrainerHandler(trainerService, cfg, zapLogger)
	gymOwnerHandler := handlers.NewGymOwnerHandler(gymOwnerService, cfg, zapLogger)
	customerHandler := handlers.NewCustomerHandler(customerService, cfg, zapLogger, otpService)
	otpHandler := handlers.NewOTPHandler(otpService, zapLogger)

	// API routes
	api := app.Group("/api")

	// Swagger documentation
	api.Get("/docs/*", swagger.New(swagger.Config{
		Title:        "Gym Management System API",
		DeepLinking:  true,
		DocExpansion: "list",
	}))

	// API v1 routes
	v1 := api.Group("/v1")

	// Health check route
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// OTP routes
	otp := v1.Group("/otp")
	otp.Post("/send", otpHandler.SendOTP)
	otp.Post("/verify", otpHandler.VerifyOTP)

	// Trainer routes
	trainers := v1.Group("/trainers")
	trainers.Get("/", trainerHandler.ListTrainers)
	trainers.Get("/:id", trainerHandler.GetTrainer)
	trainers.Post("/", middleware.GymOwnerAuthMiddleware(jwtService), trainerHandler.CreateTrainer)
	trainers.Put("/:id", middleware.AuthMiddleware(jwtService), trainerHandler.UpdateTrainer)
	trainers.Delete("/:id", middleware.AuthMiddleware(jwtService), trainerHandler.DeleteTrainer)

	// Gym owner routes
	gymOwners := v1.Group("/gym-owners")
	gymOwners.Post("/register", gymOwnerHandler.Register)
	gymOwners.Post("/login", gymOwnerHandler.Login)
	gymOwners.Get("/:id", middleware.GymOwnerAuthMiddleware(jwtService), gymOwnerHandler.GetByID)
	gymOwners.Put("/:id", middleware.GymOwnerAuthMiddleware(jwtService), gymOwnerHandler.Update)
	gymOwners.Delete("/:id", gymOwnerHandler.Delete)

	// Gym owner trainer management routes
	gymOwners.Post("/trainers", middleware.GymOwnerAuthMiddleware(jwtService), gymOwnerHandler.CreateTrainerForGym)
	gymOwners.Get("/trainers", middleware.GymOwnerAuthMiddleware(jwtService), gymOwnerHandler.ListGymTrainers)

	// Gym owner customer management routes
	gymOwners.Post("/customers", middleware.GymOwnerAuthMiddleware(jwtService), gymOwnerHandler.CreateCustomerForGym)
	gymOwners.Get("/customers", middleware.GymOwnerAuthMiddleware(jwtService), gymOwnerHandler.ListGymCustomers)

	// Customer routes
	customers := v1.Group("/customers")
	customers.Get("/", customerHandler.ListCustomers)
	customers.Get("/:id", customerHandler.GetCustomer)
	customers.Get("/by-goal/:goal", customerHandler.ListCustomersByGoal)
	customers.Post("/", middleware.GymOwnerAuthMiddleware(jwtService), customerHandler.Register)
	customers.Put("/:id", middleware.AuthMiddleware(jwtService), customerHandler.UpdateCustomer)
	customers.Delete("/:id", middleware.AuthMiddleware(jwtService), customerHandler.DeleteCustomer)

	// Start server
	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
			zapLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		zapLogger.Fatal("Failed to shutdown server", zap.Error(err))
	}

	// Close database connection
	if err := database.Close(); err != nil {
		zapLogger.Error("Failed to close database connection", zap.Error(err))
	}
}

package handlers

import (
	// "strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yourname/fitness-platform/internal/config"
	"github.com/yourname/fitness-platform/internal/models"
	"github.com/yourname/fitness-platform/internal/service"
	"go.uber.org/zap"
)

// GymOwnerHandler handles HTTP requests for gym owner operations
type GymOwnerHandler struct {
	service *service.GymOwnerService
	config  *config.Config
	logger  *zap.Logger
}

// NewGymOwnerHandler creates a new GymOwnerHandler
func NewGymOwnerHandler(service *service.GymOwnerService, cfg *config.Config, logger *zap.Logger) *GymOwnerHandler {
	return &GymOwnerHandler{
		service: service,
		config:  cfg,
		logger:  logger,
	}
}

// @Summary Register a new gym owner
// @Description Register a new gym owner with the provided information
// @Tags gym-owners
// @Accept json
// @Produce json
// @Param request body models.CreateGymOwnerRequest true "Gym owner registration details"
// @Success 201 {object} models.GymOwner
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /gym-owners/register [post]
func (h *GymOwnerHandler) Register(c *fiber.Ctx) error {
	var req models.CreateGymOwnerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid request body",
			fiber.StatusBadRequest,
			err,
		))
	}

	owner, err := h.service.Register(c.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to register gym owner", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to register gym owner",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success:    true,
		Message:    "Gym owner registered successfully",
		StatusCode: fiber.StatusCreated,
		Data:       owner,
	})
}

// @Summary Login gym owner
// @Description Authenticate a gym owner and return a JWT token
// @Tags gym-owners
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /gym-owners/login [post]
func (h *GymOwnerHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid request body",
			fiber.StatusBadRequest,
			err,
		))
	}

	token, err := h.service.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.Error("Failed to login gym owner", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(models.NewErrorResponse(
			"Invalid credentials",
			fiber.StatusUnauthorized,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "Login successful",
		StatusCode: fiber.StatusOK,
		Data: fiber.Map{
			"token": token,
		},
	})
}

// @Summary Get gym owner profile
// @Description Get gym owner details by ID
// @Tags gym-owners
// @Accept json
// @Produce json
// @Param id path string true "Gym Owner ID"
// @Success 200 {object} models.GymOwner
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /gym-owners/{id} [get]
func (h *GymOwnerHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid ID format",
			fiber.StatusBadRequest,
			err,
		))
	}

	owner, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get gym owner", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse(
			"Gym owner not found",
			fiber.StatusNotFound,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "Gym owner retrieved successfully",
		StatusCode: fiber.StatusOK,
		Data:       owner,
	})
}

// @Summary Update gym owner profile
// @Description Update gym owner information
// @Tags gym-owners
// @Accept json
// @Produce json
// @Param id path string true "Gym Owner ID"
// @Param request body models.UpdateGymOwnerRequest true "Update details"
// @Success 200 {object} models.GymOwner
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /gym-owners/{id} [put]
func (h *GymOwnerHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid ID format",
			fiber.StatusBadRequest,
			err,
		))
	}

	var req models.UpdateGymOwnerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid request body",
			fiber.StatusBadRequest,
			err,
		))
	}

	owner, err := h.service.Update(c.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to update gym owner", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to update gym owner",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "Gym owner updated successfully",
		StatusCode: fiber.StatusOK,
		Data:       owner,
	})
}

// @Summary Delete a gym owner
// @Description Delete a gym owner. This endpoint is public and does not require authentication.
// @Tags gym-owners
// @Accept json
// @Produce json
// @Param id path string true "Gym Owner ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /gym-owners/{id} [delete]
func (h *GymOwnerHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	ownerID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid gym owner ID",
			fiber.StatusBadRequest,
			err,
		))
	}

	if err := h.service.Delete(c.Context(), ownerID); err != nil {
		if err == service.ErrGymOwnerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse(
				"Gym owner not found",
				fiber.StatusNotFound,
				err,
			))
		}
		h.logger.Error("Failed to delete gym owner", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to delete gym owner",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "Gym owner deleted successfully",
		StatusCode: fiber.StatusOK,
	})
}

// @Summary Verify OTP
// @Description Verify OTP for gym owner account activation
// @Tags gym-owners
// @Accept json
// @Produce json
// @Param request body models.VerifyOTPRequest true "OTP verification details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /gym-owners/verify-otp [post]
func (h *GymOwnerHandler) VerifyOTP(c *fiber.Ctx) error {
	var req models.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid request body",
			fiber.StatusBadRequest,
			err,
		))
	}

	ownerID, err := uuid.Parse(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid user ID format",
			fiber.StatusBadRequest,
			err,
		))
	}

	if err := h.service.VerifyOTP(c.Context(), ownerID, req.OTP); err != nil {
		h.logger.Error("Failed to verify OTP", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid OTP",
			fiber.StatusBadRequest,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "OTP verified successfully",
		StatusCode: fiber.StatusOK,
	})
}

// CreateTrainerForGym creates a new trainer under the gym owner
// @Summary Create trainer for gym
// @Description Create a new trainer under the gym owner's management
// @Tags gym-owners
// @Accept json
// @Produce json
// @Param request body models.CreateTrainerRequest true "Trainer details"
// @Success 201 {object} models.Trainer
// @Failure 400,401,403 {object} map[string]string
// @Security BearerAuth
// @Router /gym-owners/trainers [post]
func (h *GymOwnerHandler) CreateTrainerForGym(c *fiber.Ctx) error {
	var req models.CreateTrainerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid request body",
			fiber.StatusBadRequest,
			err,
		))
	}

	gymOwnerID := c.Locals("userID").(string)
	trainer, err := h.service.CreateTrainerForGym(c.Context(), gymOwnerID, &req)
	if err != nil {
		h.logger.Error("Failed to create trainer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to create trainer",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success:    true,
		Message:    "Trainer created successfully",
		StatusCode: fiber.StatusCreated,
		Data:       trainer,
	})
}

// CreateCustomerForGym creates a new customer under the gym owner
// @Summary Create customer for gym
// @Description Create a new customer under the gym owner's management
// @Tags gym-owners
// @Accept json
// @Produce json
// @Param request body models.CreateCustomerRequest true "Customer details"
// @Success 201 {object} models.Customer
// @Failure 400,401,403 {object} map[string]string
// @Security BearerAuth
// @Router /gym-owners/customers [post]
func (h *GymOwnerHandler) CreateCustomerForGym(c *fiber.Ctx) error {
	var req models.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid request body",
			fiber.StatusBadRequest,
			err,
		))
	}

	gymOwnerID := c.Locals("userID").(string)
	customer, err := h.service.CreateCustomerForGym(c.Context(), gymOwnerID, &req)
	if err != nil {
		h.logger.Error("Failed to create customer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to create customer",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success:    true,
		Message:    "Customer created successfully",
		StatusCode: fiber.StatusCreated,
		Data:       customer,
	})
}

// ListGymTrainers lists all trainers under the gym owner
// @Summary List gym trainers
// @Description Get a list of all trainers under the gym owner's management
// @Tags gym-owners
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} models.Trainer
// @Failure 401,403 {object} map[string]string
// @Security BearerAuth
// @Router /gym-owners/trainers [get]
func (h *GymOwnerHandler) ListGymTrainers(c *fiber.Ctx) error {
	gymOwnerID := c.Locals("userID").(string)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	trainers, total, err := h.service.ListGymTrainers(c.Context(), gymOwnerID, page, limit)
	if err != nil {
		h.logger.Error("Failed to list trainers", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to list trainers",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "Trainers retrieved successfully",
		StatusCode: fiber.StatusOK,
		Data: fiber.Map{
			"trainers": trainers,
			"total":    total,
		},
	})
}

// ListGymCustomers lists all customers under the gym owner
// @Summary List gym customers
// @Description Get a list of all customers under the gym owner's management
// @Tags gym-owners
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} models.Customer
// @Failure 401,403 {object} map[string]string
// @Security BearerAuth
// @Router /gym-owners/customers [get]
func (h *GymOwnerHandler) ListGymCustomers(c *fiber.Ctx) error {
	gymOwnerID := c.Locals("userID").(string)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	customers, total, err := h.service.ListGymCustomers(c.Context(), gymOwnerID, page, limit)
	if err != nil {
		h.logger.Error("Failed to list customers", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to list customers",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "Customers retrieved successfully",
		StatusCode: fiber.StatusOK,
		Data: fiber.Map{
			"customers": customers,
			"total":     total,
		},
	})
}

package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yourname/fitness-platform/internal/config"
	"github.com/yourname/fitness-platform/internal/models"
	"github.com/yourname/fitness-platform/internal/service"
	"go.uber.org/zap"
)

// TrainerHandler handles trainer-related HTTP requests
type TrainerHandler struct {
	trainerService *service.TrainerService
	config         *config.Config
	logger         *zap.Logger
	validator      *validator.Validate
}

// NewTrainerHandler creates a new TrainerHandler
func NewTrainerHandler(trainerService *service.TrainerService, cfg *config.Config, log *zap.Logger) *TrainerHandler {
	return &TrainerHandler{
		trainerService: trainerService,
		config:         cfg,
		logger:         log,
		validator:      validator.New(),
	}
}

// @Summary Create a new trainer
// @Description Create a new trainer. Only accessible by gym owners.
// @Tags trainers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param trainer body models.CreateTrainerRequest true "Trainer information"
// @Success 201 {object} models.Trainer
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - Not a gym owner"
// @Failure 500 {object} map[string]string
// @Router /trainers [post]
func (h *TrainerHandler) CreateTrainer(c *fiber.Ctx) error {
	var req models.CreateTrainerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid request body",
			fiber.StatusBadRequest,
			err,
		))
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success:    false,
			Message:    err.Error(),
			StatusCode: fiber.StatusBadRequest,
		})
	}

	trainer, err := h.trainerService.CreateTrainer(c.Context(), &req)
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

// @Summary Get a trainer by ID
// @Description Get detailed information about a specific trainer. This endpoint is publicly accessible.
// @Tags trainers
// @Accept json
// @Produce json
// @Param id path string true "Trainer ID"
// @Success 200 {object} models.Trainer
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trainers/{id} [get]
func (h *TrainerHandler) GetTrainer(c *fiber.Ctx) error {
	id := c.Params("id")

	trainer, err := h.trainerService.GetTrainer(c.Context(), id)
	if err != nil {
		if err == service.ErrTrainerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse(
				"Trainer not found",
				fiber.StatusNotFound,
				err,
			))
		}
		h.logger.Error("Failed to get trainer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to get trainer",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "Trainer retrieved successfully",
		StatusCode: fiber.StatusOK,
		Data:       trainer,
	})
}

// @Summary Update a trainer
// @Description Update trainer information. Only accessible by authenticated users.
// @Tags trainers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Trainer ID"
// @Param trainer body models.UpdateTrainerRequest true "Trainer information"
// @Success 200 {object} models.Trainer
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trainers/{id} [put]
func (h *TrainerHandler) UpdateTrainer(c *fiber.Ctx) error {
	id := c.Params("id")

	var req models.UpdateTrainerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid request body",
			fiber.StatusBadRequest,
			err,
		))
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success:    false,
			Message:    err.Error(),
			StatusCode: fiber.StatusBadRequest,
		})
	}

	trainer, err := h.trainerService.UpdateTrainer(c.Context(), id, &req)
	if err != nil {
		if err == service.ErrTrainerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse(
				"Trainer not found",
				fiber.StatusNotFound,
				err,
			))
		}
		h.logger.Error("Failed to update trainer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to update trainer",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "Trainer updated successfully",
		StatusCode: fiber.StatusOK,
		Data:       trainer,
	})
}

// @Summary Delete a trainer
// @Description Soft delete a trainer. Only accessible by authenticated users.
// @Tags trainers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Trainer ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trainers/{id} [delete]
func (h *TrainerHandler) DeleteTrainer(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get authenticated user ID from context
	userID := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Success:    false,
			Message:    "Unauthorized",
			StatusCode: fiber.StatusUnauthorized,
		})
	}

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success:    false,
			Message:    "Invalid trainer ID",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// Get the trainer to check ownership
	trainer, err := h.trainerService.GetTrainer(c.Context(), id)
	if err != nil {
		if err == service.ErrTrainerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse(
				"Trainer not found",
				fiber.StatusNotFound,
				err,
			))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to get trainer",
			fiber.StatusInternalServerError,
			err,
		))
	}

	// Only allow deletion if the user is the trainer themselves
	if userID != trainer.ID.String() {
		return c.Status(fiber.StatusForbidden).JSON(models.Response{
			Success:    false,
			Message:    "You can only delete your own account",
			StatusCode: fiber.StatusForbidden,
		})
	}

	if err := h.trainerService.DeleteTrainer(c.Context(), id); err != nil {
		h.logger.Error("Failed to delete trainer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to delete trainer",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "Trainer deleted successfully",
		StatusCode: fiber.StatusOK,
	})
}

// @Summary List trainers
// @Description Get a paginated list of all trainers. This endpoint is publicly accessible.
// @Tags trainers
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {object} models.TrainerListResponse
// @Failure 500 {object} map[string]string
// @Router /trainers [get]
func (h *TrainerHandler) ListTrainers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	trainers, total, err := h.trainerService.ListTrainers(c.Context(), page, limit)
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

// @Summary Verify trainer OTP
// @Description Verify OTP for trainer activation
// @Tags trainers
// @Accept json
// @Produce json
// @Param request body models.VerifyOTPRequest true "OTP verification request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trainers/verify-otp [post]
func (h *TrainerHandler) VerifyOTP(c *fiber.Ctx) error {
	var req models.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid request body",
			fiber.StatusBadRequest,
			err,
		))
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success:    false,
			Message:    err.Error(),
			StatusCode: fiber.StatusBadRequest,
		})
	}

	trainerID, err := uuid.Parse(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success:    false,
			Message:    "Invalid trainer ID",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	err = h.trainerService.VerifyOTP(c.Context(), trainerID, req.OTP)
	if err != nil {
		h.logger.Error("Failed to verify trainer OTP", zap.Error(err))
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
		Data: fiber.Map{
			"trainer_id": req.UserID,
			"is_active":  true,
		},
	})
}

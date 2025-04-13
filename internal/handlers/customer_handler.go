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

// CustomerHandler handles customer-related HTTP requests
type CustomerHandler struct {
	customerService *service.CustomerService
	config          *config.Config
	logger          *zap.Logger
	validator       *validator.Validate
	otpService      *service.OTPService
}

// NewCustomerHandler creates a new CustomerHandler
func NewCustomerHandler(customerService *service.CustomerService, cfg *config.Config, log *zap.Logger, otpService *service.OTPService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
		config:          cfg,
		logger:          log,
		validator:       validator.New(),
		otpService:      otpService,
	}
}

// @Summary Register a new customer
// @Description Register a new customer. Only accessible by gym owners.
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param customer body models.CreateCustomerRequest true "Customer information"
// @Success 201 {object} models.Customer
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - Not a gym owner"
// @Failure 500 {object} map[string]string
// @Router /customers [post]
func (h *CustomerHandler) Register(c *fiber.Ctx) error {
	var req models.CreateCustomerRequest
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

	customer, err := h.customerService.Register(c.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to register customer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to register customer",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success:    true,
		Message:    "Customer registered successfully",
		StatusCode: fiber.StatusCreated,
		Data:       customer,
	})
}

// @Summary Get a customer by ID
// @Description Get customer details by ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} models.Customer
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomer(c *fiber.Ctx) error {
	id := c.Params("id")

	customerID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success:    false,
			Message:    "Invalid customer ID",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	customer, err := h.customerService.GetByID(c.Context(), customerID)
	if err != nil {
		if err == service.ErrCustomerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse(
				"Customer not found",
				fiber.StatusNotFound,
				err,
			))
		}
		h.logger.Error("Failed to get customer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to get customer",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "Customer retrieved successfully",
		StatusCode: fiber.StatusOK,
		Data:       customer,
	})
}

// @Summary Update a customer
// @Description Update customer details. Only accessible by authenticated users.
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Customer ID"
// @Param customer body models.UpdateCustomerRequest true "Customer information"
// @Success 200 {object} models.Customer
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
	id := c.Params("id")

	customerID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success:    false,
			Message:    "Invalid customer ID",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	var req models.UpdateCustomerRequest
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

	customer, err := h.customerService.Update(c.Context(), customerID, &req)
	if err != nil {
		if err == service.ErrCustomerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse(
				"Customer not found",
				fiber.StatusNotFound,
				err,
			))
		}
		h.logger.Error("Failed to update customer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to update customer",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "Customer updated successfully",
		StatusCode: fiber.StatusOK,
		Data:       customer,
	})
}

// @Summary Delete a customer
// @Description Soft delete a customer. Only accessible by authenticated users.
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Customer ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *fiber.Ctx) error {
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
			Message:    "Invalid customer ID",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// Get the customer to check ownership
	customer, err := h.customerService.GetByID(c.Context(), uuid.MustParse(id))
	if err != nil {
		if err == service.ErrCustomerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.NewErrorResponse(
				"Customer not found",
				fiber.StatusNotFound,
				err,
			))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to get customer",
			fiber.StatusInternalServerError,
			err,
		))
	}

	// Only allow deletion if the user is the customer themselves
	if userID != customer.ID.String() {
		return c.Status(fiber.StatusForbidden).JSON(models.Response{
			Success:    false,
			Message:    "You can only delete your own account",
			StatusCode: fiber.StatusForbidden,
		})
	}

	if err := h.customerService.Delete(c.Context(), uuid.MustParse(id)); err != nil {
		h.logger.Error("Failed to delete customer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to delete customer",
			fiber.StatusInternalServerError,
			err,
		))
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success:    true,
		Message:    "Customer deleted successfully",
		StatusCode: fiber.StatusOK,
	})
}

// @Summary List customers
// @Description Get a list of customers with pagination
// @Tags customers
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {array} models.Customer
// @Failure 500 {object} map[string]string
// @Router /customers [get]
func (h *CustomerHandler) ListCustomers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	customers, total, err := h.customerService.List(c.Context(), page, limit)
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
			"page":      page,
			"limit":     limit,
		},
	})
}

// @Summary List customers by fitness goal
// @Description Get a list of customers with a specific fitness goal
// @Tags customers
// @Accept json
// @Produce json
// @Param goal path string true "Fitness goal"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {array} models.Customer
// @Failure 500 {object} map[string]string
// @Router /customers/by-goal/{goal} [get]
func (h *CustomerHandler) ListCustomersByGoal(c *fiber.Ctx) error {
	goal := c.Params("goal")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	customers, total, err := h.customerService.ListByFitnessGoal(c.Context(), goal, page, limit)
	if err != nil {
		h.logger.Error("Failed to list customers by goal", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
			"Failed to list customers by goal",
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
			"page":      page,
			"limit":     limit,
			"goal":      goal,
		},
	})
}

// VerifyOTP godoc
// @Summary Verify customer OTP
// @Description Verify OTP for customer account activation
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param request body models.VerifyOTPRequest true "OTP verification request"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /api/v1/customers/{id}/verify-otp [post]
// @Security BearerAuth
func (h *CustomerHandler) VerifyOTP(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success:    false,
			Message:    "Invalid customer ID",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	var req models.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
			"Invalid request body",
			fiber.StatusBadRequest,
			err,
		))
	}

	// Verify OTP using OTP service
	verified, err := h.otpService.VerifyOTP(c.Context(), userID.String(), models.OTPTargetCustomer, req.OTP)
	if err != nil {
		switch err {
		case models.ErrInvalidOTP:
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
				"Invalid OTP",
				fiber.StatusBadRequest,
				err,
			))
		case models.ErrOTPExpired:
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(
				"OTP has expired",
				fiber.StatusBadRequest,
				err,
			))
		default:
			h.logger.Error("failed to verify OTP",
				zap.String("user_id", userID.String()),
				zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(
				"Failed to verify OTP",
				fiber.StatusInternalServerError,
				err,
			))
		}
	}

	if !verified {
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

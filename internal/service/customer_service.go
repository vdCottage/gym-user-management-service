package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourname/fitness-platform/internal/cache"
	"github.com/yourname/fitness-platform/internal/config"
	"github.com/yourname/fitness-platform/internal/models"
	"github.com/yourname/fitness-platform/internal/repository"
	"go.uber.org/zap"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrInvalidCustomer  = errors.New("invalid customer credentials")
)

// CustomerService handles business logic for customers
type CustomerService struct {
	customerRepo *repository.CustomerRepository
	redisClient  cache.RedisClient
	config       *config.Config
	logger       *zap.Logger
}

// NewCustomerService creates a new CustomerService
func NewCustomerService(customerRepo *repository.CustomerRepository, redisClient cache.RedisClient, config *config.Config, logger *zap.Logger) (*CustomerService, error) {
	if customerRepo == nil || redisClient == nil || config == nil || logger == nil {
		return nil, errors.New("invalid dependencies")
	}
	return &CustomerService{
		customerRepo: customerRepo,
		redisClient:  redisClient,
		config:       config,
		logger:       logger,
	}, nil
}

// Register registers a new customer
func (s *CustomerService) Register(ctx context.Context, req *models.CreateCustomerRequest) (*models.Customer, error) {
	// Check if customer with same email or phone exists
	if req.Email != "" {
		if _, err := s.customerRepo.GetByEmail(ctx, req.Email); err == nil {
			s.logger.Error("email already exists", zap.String("email", req.Email))
			return nil, fmt.Errorf("email already exists")
		}
	}

	if _, err := s.customerRepo.GetByPhone(ctx, req.Phone); err == nil {
		s.logger.Error("phone already exists", zap.String("phone", req.Phone))
		return nil, fmt.Errorf("phone already exists")
	}

	customer := &models.Customer{
		FullName:         req.FullName,
		Email:            req.Email,
		Phone:            req.Phone,
		Gender:           string(req.Gender),
		Age:              req.Age,
		Height:           float64(req.Height),
		Weight:           float64(req.Weight),
		HealthConditions: req.HealthConditions,
		FitnessGoals:     req.FitnessGoals,
		ProfileURL:       req.ProfileURL,
		IsActive:         false,
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		s.logger.Error("failed to create customer", zap.Error(err))
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	s.logger.Info("customer registered successfully", zap.String("customer_id", customer.ID.String()))
	return customer, nil
}

// GetByID retrieves a customer by ID
func (s *CustomerService) GetByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get customer", zap.Error(err), zap.String("customer_id", id.String()))
		return nil, ErrCustomerNotFound
	}
	return customer, nil
}

// Update updates a customer's information
func (s *CustomerService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateCustomerRequest) (*models.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get customer for update", zap.Error(err), zap.String("customer_id", id.String()))
		return nil, ErrCustomerNotFound
	}

	// Check if new email or phone already exists
	if req.Email != "" && req.Email != customer.Email {
		if _, err := s.customerRepo.GetByEmail(ctx, req.Email); err == nil {
			s.logger.Error("email already exists", zap.String("email", req.Email))
			return nil, fmt.Errorf("email already exists")
		}
		customer.Email = req.Email
	}

	if req.Phone != "" && req.Phone != customer.Phone {
		if _, err := s.customerRepo.GetByPhone(ctx, req.Phone); err == nil {
			s.logger.Error("phone already exists", zap.String("phone", req.Phone))
			return nil, fmt.Errorf("phone already exists")
		}
		customer.Phone = req.Phone
	}

	// Update other fields if provided
	if req.FullName != "" {
		customer.FullName = req.FullName
	}
	if req.Gender != "" {
		customer.Gender = string(req.Gender)
	}
	if req.Age > 0 {
		customer.Age = req.Age
	}
	if req.Height > 0 {
		customer.Height = float64(req.Height)
	}
	if req.Weight > 0 {
		customer.Weight = float64(req.Weight)
	}
	if req.HealthConditions != "" {
		customer.HealthConditions = req.HealthConditions
	}
	if req.FitnessGoals != "" {
		customer.FitnessGoals = req.FitnessGoals
	}
	if req.ProfileURL != "" {
		customer.ProfileURL = req.ProfileURL
	}

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		s.logger.Error("failed to update customer", zap.Error(err), zap.String("customer_id", id.String()))
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	s.logger.Info("customer updated successfully", zap.String("customer_id", id.String()))
	return customer, nil
}

// Delete deletes a customer
func (s *CustomerService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.customerRepo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete customer", zap.Error(err), zap.String("customer_id", id.String()))
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	s.logger.Info("customer deleted successfully", zap.String("customer_id", id.String()))
	return nil
}

// List retrieves a list of customers with pagination
func (s *CustomerService) List(ctx context.Context, page, limit int) ([]*models.Customer, int64, error) {
	customers, total, err := s.customerRepo.List(ctx, page, limit)
	if err != nil {
		s.logger.Error("failed to list customers", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}
	return customers, total, nil
}

// ListByFitnessGoal retrieves a list of customers by fitness goal with pagination
func (s *CustomerService) ListByFitnessGoal(ctx context.Context, goal string, page, limit int) ([]*models.Customer, int64, error) {
	customers, total, err := s.customerRepo.ListByFitnessGoal(ctx, goal, page, limit)
	if err != nil {
		s.logger.Error("failed to list customers by fitness goal", zap.Error(err), zap.String("goal", goal))
		return nil, 0, fmt.Errorf("failed to list customers by fitness goal: %w", err)
	}
	return customers, total, nil
}

// Activate activates a customer account
func (s *CustomerService) Activate(ctx context.Context, id uuid.UUID) error {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get customer for activation", zap.Error(err), zap.String("customer_id", id.String()))
		return ErrCustomerNotFound
	}

	customer.IsActive = true
	if err := s.customerRepo.Update(ctx, customer); err != nil {
		s.logger.Error("failed to activate customer", zap.Error(err), zap.String("customer_id", id.String()))
		return fmt.Errorf("failed to activate customer: %w", err)
	}

	s.logger.Info("customer activated successfully", zap.String("customer_id", id.String()))
	return nil
}

// Deactivate deactivates a customer account
func (s *CustomerService) Deactivate(ctx context.Context, id uuid.UUID) error {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get customer for deactivation", zap.Error(err), zap.String("customer_id", id.String()))
		return ErrCustomerNotFound
	}

	customer.IsActive = false
	if err := s.customerRepo.Update(ctx, customer); err != nil {
		s.logger.Error("failed to deactivate customer", zap.Error(err), zap.String("customer_id", id.String()))
		return fmt.Errorf("failed to deactivate customer: %w", err)
	}

	s.logger.Info("customer deactivated successfully", zap.String("customer_id", id.String()))
	return nil
}

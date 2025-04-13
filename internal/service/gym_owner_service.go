package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourname/fitness-platform/internal/auth"
	"github.com/yourname/fitness-platform/internal/cache"
	"github.com/yourname/fitness-platform/internal/config"
	"github.com/yourname/fitness-platform/internal/models"
	"github.com/yourname/fitness-platform/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrGymOwnerNotFound = errors.New("gym owner not found")
	ErrInvalidGymOwner  = errors.New("invalid gym owner credentials")
)

// GymOwnerService handles business logic for gym owners
type GymOwnerService struct {
	gymOwnerRepo *repository.GymOwnerRepository
	trainerRepo  *repository.TrainerRepository
	customerRepo *repository.CustomerRepository
	redisClient  cache.RedisClient
	config       *config.Config
	logger       *zap.Logger
}

// NewGymOwnerService creates a new GymOwnerService
func NewGymOwnerService(
	gymOwnerRepo *repository.GymOwnerRepository,
	trainerRepo *repository.TrainerRepository,
	customerRepo *repository.CustomerRepository,
	redisClient cache.RedisClient,
	config *config.Config,
	logger *zap.Logger,
) (*GymOwnerService, error) {
	if gymOwnerRepo == nil || trainerRepo == nil || customerRepo == nil || redisClient == nil || config == nil || logger == nil {
		return nil, errors.New("invalid dependencies")
	}
	return &GymOwnerService{
		gymOwnerRepo: gymOwnerRepo,
		trainerRepo:  trainerRepo,
		customerRepo: customerRepo,
		redisClient:  redisClient,
		config:       config,
		logger:       logger,
	}, nil
}

// Register registers a new gym owner
func (s *GymOwnerService) Register(ctx context.Context, req *models.CreateGymOwnerRequest) (*models.GymOwner, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	owner := &models.GymOwner{
		FullName:              req.FullName,
		Email:                 req.Email,
		Phone:                 req.Phone,
		Password:              string(hashedPassword),
		GymName:               req.GymName,
		GymRegistrationNumber: req.GymRegistrationNumber,
		// IsActive:              false,
	}

	if err := s.gymOwnerRepo.Create(ctx, owner); err != nil {
		s.logger.Error("failed to create gym owner", zap.Error(err))
		return nil, fmt.Errorf("failed to create gym owner: %w", err)
	}

	s.logger.Info("gym owner registered successfully", zap.String("owner_id", owner.ID.String()))
	return owner, nil
}

// Login authenticates a gym owner
func (s *GymOwnerService) Login(ctx context.Context, email, password string) (string, error) {
	owner, err := s.gymOwnerRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Warn("login failed: gym owner not found", zap.String("email", email))
		return "", ErrInvalidGymOwner
	}

	if err := bcrypt.CompareHashAndPassword([]byte(owner.Password), []byte(password)); err != nil {
		s.logger.Warn("login failed: invalid password", zap.String("email", email))
		return "", ErrInvalidGymOwner
	}

	token, err := auth.GenerateJWTToken(owner.ID, "gym_owner", s.config)
	if err != nil {
		s.logger.Error("failed to generate token", zap.Error(err))
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info("gym owner logged in successfully", zap.String("owner_id", owner.ID.String()))
	return token, nil
}

// GetByID retrieves a gym owner by ID
func (s *GymOwnerService) GetByID(ctx context.Context, id uuid.UUID) (*models.GymOwner, error) {
	owner, err := s.gymOwnerRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get gym owner", zap.Error(err), zap.String("owner_id", id.String()))
		return nil, ErrGymOwnerNotFound
	}
	return owner, nil
}

// Update updates a gym owner's information
func (s *GymOwnerService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateGymOwnerRequest) (*models.GymOwner, error) {
	owner, err := s.gymOwnerRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get gym owner for update", zap.Error(err), zap.String("owner_id", id.String()))
		return nil, ErrGymOwnerNotFound
	}

	if req.FullName != "" {
		owner.FullName = req.FullName
	}
	if req.Phone != "" {
		owner.Phone = req.Phone
	}
	if req.GymName != "" {
		owner.GymName = req.GymName
	}

	if err := s.gymOwnerRepo.Update(ctx, owner); err != nil {
		s.logger.Error("failed to update gym owner", zap.Error(err), zap.String("owner_id", id.String()))
		return nil, fmt.Errorf("failed to update gym owner: %w", err)
	}

	s.logger.Info("gym owner updated successfully", zap.String("owner_id", id.String()))
	return owner, nil
}

// Delete deletes a gym owner
func (s *GymOwnerService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.gymOwnerRepo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete gym owner", zap.Error(err), zap.String("owner_id", id.String()))
		return fmt.Errorf("failed to delete gym owner: %w", err)
	}

	s.logger.Info("gym owner deleted successfully", zap.String("owner_id", id.String()))
	return nil
}

// List retrieves a list of gym owners with pagination
func (s *GymOwnerService) List(ctx context.Context, page, limit int) ([]*models.GymOwner, int64, error) {
	owners, total, err := s.gymOwnerRepo.List(ctx, page, limit)
	if err != nil {
		s.logger.Error("failed to list gym owners", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list gym owners: %w", err)
	}
	return owners, total, nil
}

// VerifyOTP verifies the OTP and activates the gym owner
func (s *GymOwnerService) VerifyOTP(ctx context.Context, ownerID uuid.UUID, otp string) error {
	// Get stored OTP from Redis
	key := fmt.Sprintf("otp:gym_owner:%s", ownerID.String())
	storedOTP, err := s.redisClient.Get(ctx, key)
	if err != nil {
		s.logger.Error("failed to get OTP", zap.Error(err))
		return ErrInvalidOTP
	}

	if storedOTP != otp {
		return ErrInvalidOTP
	}

	// // Get owner and update status
	// owner, err := s.gymOwnerRepo.GetByID(ctx, ownerID)
	// if err != nil {
	// 	s.logger.Error("failed to get gym owner", zap.Error(err))
	// 	return ErrGymOwnerNotFound
	// }

	// owner.IsActive = true
	// if err := s.gymOwnerRepo.Update(ctx, owner); err != nil {
	// 	s.logger.Error("failed to update gym owner status", zap.Error(err))
	// 	return fmt.Errorf("failed to update gym owner status: %w", err)
	// }

	// Delete OTP from Redis
	if err := s.redisClient.Delete(ctx, key); err != nil {
		s.logger.Error("failed to delete OTP", zap.Error(err))
	}

	s.logger.Info("gym owner activated successfully", zap.String("owner_id", ownerID.String()))
	return nil
}

// Activate activates a gym owner account
func (s *GymOwnerService) Activate(ctx context.Context, id uuid.UUID) error {
	owner, err := s.gymOwnerRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get gym owner for activation", zap.Error(err), zap.String("owner_id", id.String()))
		return ErrGymOwnerNotFound
	}

	owner.IsActive = true
	if err := s.gymOwnerRepo.Update(ctx, owner); err != nil {
		s.logger.Error("failed to activate gym owner", zap.Error(err), zap.String("owner_id", id.String()))
		return fmt.Errorf("failed to activate gym owner: %w", err)
	}

	s.logger.Info("gym owner activated successfully", zap.String("owner_id", id.String()))
	return nil
}

// Deactivate deactivates a gym owner account
func (s *GymOwnerService) Deactivate(ctx context.Context, id uuid.UUID) error {
	owner, err := s.gymOwnerRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get gym owner for deactivation", zap.Error(err), zap.String("owner_id", id.String()))
		return ErrGymOwnerNotFound
	}

	owner.IsActive = false
	if err := s.gymOwnerRepo.Update(ctx, owner); err != nil {
		s.logger.Error("failed to deactivate gym owner", zap.Error(err), zap.String("owner_id", id.String()))
		return fmt.Errorf("failed to deactivate gym owner: %w", err)
	}

	s.logger.Info("gym owner deactivated successfully", zap.String("owner_id", id.String()))
	return nil
}

// CreateTrainerForGym creates a new trainer under the gym owner
func (s *GymOwnerService) CreateTrainerForGym(ctx context.Context, gymOwnerID string, req *models.CreateTrainerRequest) (*models.Trainer, error) {
	ownerUUID, err := uuid.Parse(gymOwnerID)
	if err != nil {
		return nil, fmt.Errorf("invalid gym owner ID: %w", err)
	}

	// Verify gym owner exists
	owner, err := s.gymOwnerRepo.GetByID(ctx, ownerUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gym owner: %w", err)
	}
	if owner == nil {
		return nil, ErrGymOwnerNotFound
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create trainer with gym owner reference
	trainer := &models.Trainer{
		ID:             uuid.New(),
		Email:          req.Email,
		Password:       string(hashedPassword),
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Phone:          req.Phone,
		GymOwnerID:     ownerUUID,
		Role:           models.RoleTrainer,
		IsActive:       true,
		Specialization: req.Specialization,
		Experience:     req.Experience,
		Bio:            req.Bio,
		IsAvailable:    true,
		Rating:         0,
	}

	if err := s.trainerRepo.Create(ctx, trainer); err != nil {
		return nil, fmt.Errorf("failed to create trainer: %w", err)
	}

	s.logger.Info("trainer created successfully",
		zap.String("trainer_id", trainer.ID.String()),
		zap.String("gym_owner_id", gymOwnerID))
	return trainer, nil
}

// CreateCustomerForGym creates a new customer under the gym owner
func (s *GymOwnerService) CreateCustomerForGym(ctx context.Context, gymOwnerID string, req *models.CreateCustomerRequest) (*models.Customer, error) {
	ownerUUID, err := uuid.Parse(gymOwnerID)
	if err != nil {
		return nil, fmt.Errorf("invalid gym owner ID: %w", err)
	}

	// Verify gym owner exists
	owner, err := s.gymOwnerRepo.GetByID(ctx, ownerUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gym owner: %w", err)
	}
	if owner == nil {
		return nil, ErrGymOwnerNotFound
	}

	// Create customer with gym owner reference
	customer := &models.Customer{
		ID:               uuid.New(),
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
		GymOwnerID:       ownerUUID,
		Role:             models.RoleCustomer,
		IsActive:         true,
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	s.logger.Info("customer created successfully",
		zap.String("customer_id", customer.ID.String()),
		zap.String("gym_owner_id", gymOwnerID))
	return customer, nil
}

// ListGymTrainers lists all trainers under the gym owner
func (s *GymOwnerService) ListGymTrainers(ctx context.Context, gymOwnerID string, page, limit int) ([]*models.Trainer, int64, error) {
	ownerUUID, err := uuid.Parse(gymOwnerID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid gym owner ID: %w", err)
	}

	// Verify gym owner exists
	owner, err := s.gymOwnerRepo.GetByID(ctx, ownerUUID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get gym owner: %w", err)
	}
	if owner == nil {
		return nil, 0, ErrGymOwnerNotFound
	}

	return s.trainerRepo.ListByGymOwner(ctx, ownerUUID, page, limit)
}

// ListGymCustomers lists all customers under the gym owner
func (s *GymOwnerService) ListGymCustomers(ctx context.Context, gymOwnerID string, page, limit int) ([]*models.Customer, int64, error) {
	ownerUUID, err := uuid.Parse(gymOwnerID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid gym owner ID: %w", err)
	}

	// Verify gym owner exists
	owner, err := s.gymOwnerRepo.GetByID(ctx, ownerUUID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get gym owner: %w", err)
	}
	if owner == nil {
		return nil, 0, ErrGymOwnerNotFound
	}

	return s.customerRepo.ListByGymOwner(ctx, ownerUUID, page, limit)
}

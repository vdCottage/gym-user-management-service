package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourname/fitness-platform/internal/config"
	"github.com/yourname/fitness-platform/internal/models"
	"github.com/yourname/fitness-platform/internal/repository"
	"github.com/yourname/fitness-platform/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrTrainerNotFound = errors.New("trainer not found")
	ErrInvalidInput    = errors.New("invalid input")
)

// TrainerService handles trainer-related business logic
type TrainerService struct {
	trainerRepo *repository.TrainerRepository
	config      *config.Config
	logger      *logger.Logger
	otpService  *OTPService
}

// NewTrainerService creates a new TrainerService
func NewTrainerService(trainerRepo *repository.TrainerRepository, cfg *config.Config, log *logger.Logger, otpService *OTPService) *TrainerService {
	return &TrainerService{
		trainerRepo: trainerRepo,
		config:      cfg,
		logger:      log,
		otpService:  otpService,
	}
}

// CreateTrainer creates a new trainer
func (s *TrainerService) CreateTrainer(ctx context.Context, req *models.CreateTrainerRequest) (*models.Trainer, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create trainer
	trainer := &models.Trainer{
		ID:             uuid.New(),
		Email:          req.Email,
		Password:       string(hashedPassword),
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Phone:          req.Phone,
		IsActive:       true,
		Specialization: req.Specialization,
		Experience:     req.Experience,
		Bio:            req.Bio,
		IsAvailable:    true,
		Rating:         0,
	}

	// Create trainer using repository
	if err := s.trainerRepo.Create(ctx, trainer); err != nil {
		s.logger.Error("failed to create trainer", zap.Error(err))
		return nil, fmt.Errorf("failed to create trainer: %w", err)
	}

	s.logger.Info("trainer created successfully",
		zap.String("trainer_id", trainer.ID.String()))
	return trainer, nil
}

// GetTrainer retrieves a trainer by ID
func (s *TrainerService) GetTrainer(ctx context.Context, id string) (*models.Trainer, error) {
	trainerID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidInput
	}

	trainer, err := s.trainerRepo.GetByID(ctx, trainerID)
	if err != nil {
		return nil, err
	}

	if trainer == nil {
		return nil, ErrTrainerNotFound
	}

	return trainer, nil
}

// UpdateTrainer updates a trainer's information
func (s *TrainerService) UpdateTrainer(ctx context.Context, id string, req *models.UpdateTrainerRequest) (*models.Trainer, error) {
	trainerID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidInput
	}

	trainer, err := s.trainerRepo.GetByID(ctx, trainerID)
	if err != nil {
		return nil, err
	}

	if trainer == nil {
		return nil, ErrTrainerNotFound
	}

	// Update fields
	if req.FirstName != "" {
		trainer.FirstName = req.FirstName
	}
	if req.LastName != "" {
		trainer.LastName = req.LastName
	}
	if req.Phone != "" {
		trainer.Phone = req.Phone
	}
	if req.Specialization != "" {
		trainer.Specialization = req.Specialization
	}
	if req.Experience > 0 {
		trainer.Experience = req.Experience
	}
	if req.Bio != "" {
		trainer.Bio = req.Bio
	}
	if req.IsAvailable != nil {
		trainer.IsAvailable = *req.IsAvailable
	}

	if err := s.trainerRepo.Update(ctx, trainer); err != nil {
		return nil, err
	}

	return trainer, nil
}

// DeleteTrainer deletes a trainer
func (s *TrainerService) DeleteTrainer(ctx context.Context, id string) error {
	trainerID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidInput
	}

	trainer, err := s.trainerRepo.GetByID(ctx, trainerID)
	if err != nil {
		return err
	}

	if trainer == nil {
		return ErrTrainerNotFound
	}

	return s.trainerRepo.Delete(ctx, trainerID)
}

// ListTrainers retrieves a list of trainers with pagination
func (s *TrainerService) ListTrainers(ctx context.Context, page, limit int) ([]*models.Trainer, int64, error) {
	return s.trainerRepo.List(ctx, page, limit)
}

// VerifyOTP verifies the OTP for a trainer
func (s *TrainerService) VerifyOTP(ctx context.Context, trainerID uuid.UUID, otp string) error {
	// Verify OTP using the OTP service
	verified, err := s.otpService.VerifyOTP(ctx, trainerID.String(), models.OTPTargetTrainer, otp)
	if err != nil {
		s.logger.Error("failed to verify OTP", zap.Error(err), zap.String("trainer_id", trainerID.String()))
		return fmt.Errorf("failed to verify OTP: %w", err)
	}
	if !verified {
		return fmt.Errorf("invalid OTP")
	}

	// Activate the trainer account
	trainer, err := s.trainerRepo.GetByID(ctx, trainerID)
	if err != nil {
		s.logger.Error("failed to get trainer for activation", zap.Error(err), zap.String("trainer_id", trainerID.String()))
		return ErrTrainerNotFound
	}

	trainer.IsActive = true
	if err := s.trainerRepo.Update(ctx, trainer); err != nil {
		s.logger.Error("failed to activate trainer", zap.Error(err), zap.String("trainer_id", trainerID.String()))
		return fmt.Errorf("failed to activate trainer: %w", err)
	}

	s.logger.Info("trainer OTP verified and account activated", zap.String("trainer_id", trainerID.String()))
	return nil
}

// Activate activates a trainer account
func (s *TrainerService) Activate(ctx context.Context, id uuid.UUID) error {
	trainer, err := s.trainerRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get trainer for activation", zap.Error(err), zap.String("trainer_id", id.String()))
		return ErrTrainerNotFound
	}

	trainer.IsActive = true
	if err := s.trainerRepo.Update(ctx, trainer); err != nil {
		s.logger.Error("failed to activate trainer", zap.Error(err), zap.String("trainer_id", id.String()))
		return fmt.Errorf("failed to activate trainer: %w", err)
	}

	s.logger.Info("trainer activated successfully", zap.String("trainer_id", id.String()))
	return nil
}

// Deactivate deactivates a trainer account
func (s *TrainerService) Deactivate(ctx context.Context, id uuid.UUID) error {
	trainer, err := s.trainerRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get trainer for deactivation", zap.Error(err), zap.String("trainer_id", id.String()))
		return ErrTrainerNotFound
	}

	trainer.IsActive = false
	if err := s.trainerRepo.Update(ctx, trainer); err != nil {
		s.logger.Error("failed to deactivate trainer", zap.Error(err), zap.String("trainer_id", id.String()))
		return fmt.Errorf("failed to deactivate trainer: %w", err)
	}

	s.logger.Info("trainer deactivated successfully", zap.String("trainer_id", id.String()))
	return nil
}

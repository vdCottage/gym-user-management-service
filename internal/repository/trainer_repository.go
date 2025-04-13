package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/yourname/fitness-platform/internal/models"
	"gorm.io/gorm"
)

var (
	ErrTrainerNotFound = errors.New("trainer not found")
)

// TrainerRepository handles database operations for trainers
type TrainerRepository struct {
	DB *gorm.DB
}

// NewTrainerRepository creates a new TrainerRepository
func NewTrainerRepository(db *gorm.DB) *TrainerRepository {
	return &TrainerRepository{
		DB: db,
	}
}

// Create creates a new trainer
func (r *TrainerRepository) Create(ctx context.Context, trainer *models.Trainer) error {
	return r.DB.WithContext(ctx).Create(trainer).Error
}

// GetByID retrieves a trainer by ID
func (r *TrainerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Trainer, error) {
	var trainer models.Trainer
	if err := r.DB.WithContext(ctx).First(&trainer, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTrainerNotFound
		}
		return nil, err
	}
	return &trainer, nil
}

// GetByEmail retrieves a trainer by email
func (r *TrainerRepository) GetByEmail(ctx context.Context, email string) (*models.Trainer, error) {
	var trainer models.Trainer
	if err := r.DB.WithContext(ctx).First(&trainer, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTrainerNotFound
		}
		return nil, err
	}
	return &trainer, nil
}

// GetByPhone retrieves a trainer by phone
func (r *TrainerRepository) GetByPhone(ctx context.Context, phone string) (*models.Trainer, error) {
	var trainer models.Trainer
	if err := r.DB.WithContext(ctx).First(&trainer, "phone = ?", phone).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTrainerNotFound
		}
		return nil, err
	}
	return &trainer, nil
}

// Update updates a trainer
func (r *TrainerRepository) Update(ctx context.Context, trainer *models.Trainer) error {
	return r.DB.WithContext(ctx).Save(trainer).Error
}

// Delete soft deletes a trainer
func (r *TrainerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&models.Trainer{}, "id = ?", id).Error
}

// List retrieves a list of trainers with pagination
func (r *TrainerRepository) List(ctx context.Context, page, limit int) ([]*models.Trainer, int64, error) {
	var trainers []*models.Trainer
	var total int64

	offset := (page - 1) * limit

	if err := r.DB.WithContext(ctx).Model(&models.Trainer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.WithContext(ctx).Offset(offset).Limit(limit).Find(&trainers).Error; err != nil {
		return nil, 0, err
	}

	return trainers, total, nil
}

// ListBySpecialization retrieves a list of trainers by specialization with pagination
func (r *TrainerRepository) ListBySpecialization(ctx context.Context, specialization string, page, limit int) ([]*models.Trainer, int64, error) {
	var trainers []*models.Trainer
	var total int64

	offset := (page - 1) * limit

	if err := r.DB.WithContext(ctx).
		Where("specialization = ?", specialization).
		Model(&models.Trainer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.WithContext(ctx).
		Where("specialization = ?", specialization).
		Offset(offset).Limit(limit).
		Find(&trainers).Error; err != nil {
		return nil, 0, err
	}

	return trainers, total, nil
}

// ListByGymOwner retrieves trainers for a specific gym owner with pagination
func (r *TrainerRepository) ListByGymOwner(ctx context.Context, gymOwnerID uuid.UUID, page, limit int) ([]*models.Trainer, int64, error) {
	var trainers []*models.Trainer
	var total int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.DB.WithContext(ctx).Model(&models.Trainer{}).Where("gym_owner_id = ?", gymOwnerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get trainers with pagination
	if err := r.DB.WithContext(ctx).
		Where("gym_owner_id = ?", gymOwnerID).
		Offset(offset).
		Limit(limit).
		Find(&trainers).Error; err != nil {
		return nil, 0, err
	}

	return trainers, total, nil
}

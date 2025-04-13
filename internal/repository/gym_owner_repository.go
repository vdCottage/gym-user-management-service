package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/yourname/fitness-platform/internal/models"
	"gorm.io/gorm"
)

// GymOwnerRepository handles database operations for gym owners
type GymOwnerRepository struct {
	db *gorm.DB
}

// NewGymOwnerRepository creates a new GymOwnerRepository
func NewGymOwnerRepository(db *gorm.DB) *GymOwnerRepository {
	return &GymOwnerRepository{db: db}
}

// Create creates a new gym owner
func (r *GymOwnerRepository) Create(ctx context.Context, owner *models.GymOwner) error {
	return r.db.WithContext(ctx).Create(owner).Error
}

// GetByID retrieves a gym owner by ID
func (r *GymOwnerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.GymOwner, error) {
	var owner models.GymOwner
	if err := r.db.WithContext(ctx).First(&owner, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &owner, nil
}

// GetByEmail retrieves a gym owner by email
func (r *GymOwnerRepository) GetByEmail(ctx context.Context, email string) (*models.GymOwner, error) {
	var owner models.GymOwner
	if err := r.db.WithContext(ctx).First(&owner, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &owner, nil
}

// GetByPhone retrieves a gym owner by phone
func (r *GymOwnerRepository) GetByPhone(ctx context.Context, phone string) (*models.GymOwner, error) {
	var owner models.GymOwner
	if err := r.db.WithContext(ctx).First(&owner, "phone = ?", phone).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &owner, nil
}

// Update updates a gym owner
func (r *GymOwnerRepository) Update(ctx context.Context, owner *models.GymOwner) error {
	return r.db.WithContext(ctx).Save(owner).Error
}

// Delete soft deletes a gym owner
func (r *GymOwnerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.GymOwner{}, "id = ?", id).Error
}

// List retrieves a list of gym owners with pagination
func (r *GymOwnerRepository) List(ctx context.Context, page, limit int) ([]*models.GymOwner, int64, error) {
	var owners []*models.GymOwner
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.GymOwner{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&owners).Error; err != nil {
		return nil, 0, err
	}

	return owners, total, nil
}

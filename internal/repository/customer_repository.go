package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/yourname/fitness-platform/internal/models"
	"gorm.io/gorm"
)

// CustomerRepository handles database operations for customers
type CustomerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new CustomerRepository
func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Create creates a new customer
func (r *CustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

// GetByID retrieves a customer by ID
func (r *CustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.WithContext(ctx).First(&customer, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &customer, nil
}

// GetByEmail retrieves a customer by email
func (r *CustomerRepository) GetByEmail(ctx context.Context, email string) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.WithContext(ctx).First(&customer, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &customer, nil
}

// GetByPhone retrieves a customer by phone
func (r *CustomerRepository) GetByPhone(ctx context.Context, phone string) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.WithContext(ctx).First(&customer, "phone = ?", phone).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &customer, nil
}

// Update updates a customer
func (r *CustomerRepository) Update(ctx context.Context, customer *models.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

// Delete soft deletes a customer
func (r *CustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Customer{}, "id = ?", id).Error
}

// List retrieves a list of customers with pagination
func (r *CustomerRepository) List(ctx context.Context, page, limit int) ([]*models.Customer, int64, error) {
	var customers []*models.Customer
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.Customer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

// ListByFitnessGoal retrieves a list of customers by fitness goal with pagination
func (r *CustomerRepository) ListByFitnessGoal(ctx context.Context, goal string, page, limit int) ([]*models.Customer, int64, error) {
	var customers []*models.Customer
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).
		Where("? = ANY(fitness_goals)", goal).
		Model(&models.Customer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("? = ANY(fitness_goals)", goal).
		Offset(offset).Limit(limit).
		Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

// ListByGymOwner retrieves customers for a specific gym owner with pagination
func (r *CustomerRepository) ListByGymOwner(ctx context.Context, gymOwnerID uuid.UUID, page, limit int) ([]*models.Customer, int64, error) {
	var customers []*models.Customer
	var total int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.WithContext(ctx).Model(&models.Customer{}).Where("gym_owner_id = ?", gymOwnerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get customers with pagination
	if err := r.db.WithContext(ctx).
		Where("gym_owner_id = ?", gymOwnerID).
		Offset(offset).
		Limit(limit).
		Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

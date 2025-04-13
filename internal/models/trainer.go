package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Trainer represents a trainer in the system
// @Description Trainer information including personal details and profile
type Trainer struct {
	// @Description Unique identifier for the trainer
	// @Example trainer_123e4567-e89b-12d3-a456-426614174000
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	// @Description Email address of the trainer (must be unique)
	// @Example john.trainer@gymapp.com
	Email string `json:"email" gorm:"unique;not null"`

	// @Description Phone number of the trainer (must be unique)
	// @Example +1-555-123-4567
	Phone string `json:"phone" gorm:"unique;not null"`

	// @Description Hashed password of the trainer (not exposed in JSON)
	Password string `json:"-" gorm:"not null"`

	// @Description First name of the trainer
	// @Example John
	FirstName string `json:"first_name" gorm:"not null"`

	// @Description Last name of the trainer
	// @Example Smith
	LastName string `json:"last_name" gorm:"not null"`

	// @Description Role of the user (always 'trainer')
	// @Example trainer
	Role string `json:"role" gorm:"type:varchar(20);not null;default:'trainer'"`

	// @Description Whether the trainer account is active
	// @Example false
	IsActive bool `json:"is_active" gorm:"default:true"`

	// @Description ID of the gym owner this trainer belongs to
	// @Example 123e4567-e89b-12d3-a456-426614174000
	GymOwnerID uuid.UUID `json:"gym_owner_id" gorm:"type:uuid;not null"`

	// @Description Trainer's area of expertise
	// @Example Strength Training
	Specialization string `json:"specialization" gorm:"not null"`

	// @Description Years of experience in training
	// @Example 5
	Experience int `json:"experience" gorm:"not null;default:0"`

	// @Description Detailed biography of the trainer
	// @Example Certified personal trainer with expertise in strength training and rehabilitation.
	Bio string `json:"bio" gorm:"type:text;not null"`

	// @Description Whether the trainer is currently accepting new clients
	// @Example true
	IsAvailable bool `json:"is_available" gorm:"not null;default:true"`

	// @Description Average rating from clients (0.00 to 5.00)
	// @Example 4.75
	Rating float32 `json:"rating" gorm:"not null;default:0"`

	// @Description When the trainer was created
	// @Example 2024-03-15T14:30:00Z
	CreatedAt time.Time `json:"created_at"`

	// @Description When the trainer was last updated
	// @Example 2024-03-15T14:30:00Z
	UpdatedAt time.Time `json:"updated_at"`

	// @Description Soft delete timestamp (null if not deleted)
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// CreateTrainerRequest represents the request to create a new trainer
// @Description Request body for creating a new trainer
type CreateTrainerRequest struct {
	// @Description Email address (must be unique)
	// @Example john.trainer@gymapp.com
	Email string `json:"email" validate:"required,email"`

	// @Description Password for the account
	// @Example SecurePass123!
	Password string `json:"password" validate:"required,min=8"`

	// @Description Trainer's first name
	// @Example John
	FirstName string `json:"first_name" validate:"required"`

	// @Description Trainer's last name
	// @Example Smith
	LastName string `json:"last_name" validate:"required"`

	// @Description Phone number (must be unique)
	// @Example +1-555-123-4567
	Phone string `json:"phone" validate:"required"`

	// @Description Area of expertise
	// @Example Strength Training
	Specialization string `json:"specialization" validate:"required"`

	// @Description Years of experience
	// @Example 5
	Experience int `json:"experience" validate:"required,min=0"`

	// @Description Detailed trainer biography
	// @Example Certified personal trainer with expertise in strength training and rehabilitation.
	Bio string `json:"bio" validate:"required"`
}

// UpdateTrainerRequest represents the request to update a trainer
// @Description Request body for updating an existing trainer
type UpdateTrainerRequest struct {
	// @Description Email address (must be unique)
	// @Example john.trainer@gymapp.com
	Email string `json:"email,omitempty" validate:"omitempty,email"`

	// @Description Trainer's first name
	// @Example John
	FirstName string `json:"first_name,omitempty"`

	// @Description Trainer's last name
	// @Example Smith
	LastName string `json:"last_name,omitempty"`

	// @Description Phone number (must be unique)
	// @Example +1-555-123-4567
	Phone string `json:"phone,omitempty"`

	// @Description Area of expertise
	// @Example Strength Training
	Specialization string `json:"specialization,omitempty"`

	// @Description Years of experience
	// @Example 5
	Experience int `json:"experience,omitempty" validate:"omitempty,min=0"`

	// @Description Detailed trainer biography
	// @Example Certified personal trainer with expertise in strength training and rehabilitation.
	Bio string `json:"bio,omitempty"`

	// @Description Whether the trainer is currently accepting new clients
	// @Example true
	IsAvailable *bool `json:"is_available,omitempty"`
}

// TrainerListResponse represents the paginated response for listing trainers
// @Description Paginated response containing a list of trainers
type TrainerListResponse struct {
	// @Description List of trainers
	Trainers []*Trainer `json:"trainers"`

	// @Description Total number of trainers matching the query
	// @Example 100
	Total int64 `json:"total"`

	// @Description Current page number
	// @Example 1
	Page int `json:"page"`

	// @Description Number of items per page
	// @Example 10
	Limit int `json:"limit"`
}

// TableName specifies the database table name for Trainer
func (Trainer) TableName() string {
	return "trainers"
}

// BeforeCreate is called before creating a new trainer
func (t *Trainer) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// MarshalJSON customizes the JSON output for Trainer
func (t Trainer) MarshalJSON() ([]byte, error) {
	type Alias Trainer
	return json.Marshal(&struct {
		ID string `json:"id"`
		Alias
	}{
		ID:    fmt.Sprintf("trainer_%s", t.ID.String()),
		Alias: (Alias)(t),
	})
}

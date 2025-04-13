package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GymOwner represents a gym owner in the system
// @Description Gym owner information
type GymOwner struct {
	// @Description Unique identifier for the gym owner
	// @example 123e4567-e89b-12d3-a456-426614174000
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" validate:"required,uuid4"`

	// @Description Full name of the gym owner
	// @example John Doe
	FullName string `json:"full_name" gorm:"not null" validate:"required"`

	// @Description Email address of the gym owner
	// @example john.doe@example.com
	Email string `json:"email" gorm:"unique;not null" validate:"required,email"`

	// @Description Phone number of the gym owner
	// @example +1234567890
	Phone string `json:"phone" gorm:"unique;not null" validate:"required"`

	// @Description Hashed password of the gym owner
	Password string `json:"-" gorm:"not null" validate:"required"`

	// @Description Name of the gym
	// @example Fitness First
	GymName string `json:"gym_name" gorm:"not null" validate:"required"`

	// @Description Registration number of the gym
	// @example GYM123456
	GymRegistrationNumber string `json:"gym_registration_number" gorm:"unique;not null" validate:"required"`

	// @Description Whether the gym owner is active
	// @example false
	IsActive bool `json:"is_active" gorm:"not null;default:false"`

	// @Description Creation timestamp
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// @Description Last update timestamp
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// @Description Soft delete timestamp
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// MarshalJSON customizes the JSON response for GymOwner
func (g GymOwner) MarshalJSON() ([]byte, error) {
	type Alias GymOwner
	return json.Marshal(&struct {
		ID string `json:"id"`
		Alias
	}{
		ID:    fmt.Sprintf("owner_%s", g.ID.String()),
		Alias: (Alias)(g),
	})
}

// TableName specifies the table name for the GymOwner model
func (GymOwner) TableName() string {
	return "gym_owners"
}

// CreateGymOwnerRequest represents a request to create a new gym owner
// @Description Request to create a new gym owner
type CreateGymOwnerRequest struct {
	// @Description Full name of the gym owner
	// @example John Doe
	FullName string `json:"full_name" validate:"required,min=2,max=100"`

	// @Description Email address of the gym owner
	// @example john.doe@example.com
	Email string `json:"email" validate:"required,email"`

	// @Description Phone number of the gym owner
	// @example +1234567890
	Phone string `json:"phone" validate:"required,e164"`

	// @Description Password for the gym owner account
	// @example securePassword123
	Password string `json:"password" validate:"required,min=8,max=100"`

	// @Description Name of the gym
	// @example Fitness First
	GymName string `json:"gym_name" validate:"required,min=2,max=100"`

	// @Description Registration number of the gym
	// @example GYM123456
	GymRegistrationNumber string `json:"gym_registration_number" validate:"required,min=5,max=20"`
}

// UpdateGymOwnerRequest represents a request to update a gym owner
// @Description Request to update a gym owner
type UpdateGymOwnerRequest struct {
	// @Description Full name of the gym owner
	// @example John Doe
	FullName string `json:"full_name,omitempty" validate:"omitempty,min=2,max=100"`

	// @Description Phone number of the gym owner
	// @example +1234567890
	Phone string `json:"phone,omitempty" validate:"omitempty,e164"`

	// @Description Name of the gym
	// @example Fitness First
	GymName string `json:"gym_name,omitempty" validate:"omitempty,min=2,max=100"`
}

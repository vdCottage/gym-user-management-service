package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Gender represents the gender of a customer
type Gender string

const (
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
	GenderOther   Gender = "other"
	GenderUnknown Gender = "unknown"
)

// Customer represents a customer in the system
// @Description Customer information
type Customer struct {
	// @Description Unique identifier for the customer
	// @example 123e4567-e89b-12d3-a456-426614174000
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`

	// @Description Full name of the customer
	// @example John Doe
	FullName string `json:"full_name" gorm:"not null"`

	// @Description Email address of the customer
	// @example john.doe@example.com
	Email string `json:"email" gorm:"unique;not null"`

	// @Description Phone number of the customer
	// @example +1234567890
	Phone string `json:"phone" gorm:"unique;not null"`

	// @Description Gender of the customer
	// @example male
	Gender string `json:"gender" gorm:"not null"`

	// @Description Age of the customer
	// @example 25
	Age int `json:"age" gorm:"not null"`

	// @Description Height of the customer in centimeters
	// @example 175
	Height float64 `json:"height" gorm:"not null"`

	// @Description Weight of the customer in kilograms
	// @example 70
	Weight float64 `json:"weight" gorm:"not null"`

	// @Description Health conditions of the customer
	// @example ["Asthma", "Back Pain"]
	HealthConditions string `json:"health_conditions"`

	// @Description Fitness goals of the customer
	// @example ["Weight Loss", "Muscle Gain"]
	FitnessGoals string `json:"fitness_goals"`

	// @Description Profile URL of the customer
	// @example https://example.com/profile.jpg
	ProfileURL string `json:"profile_url" gorm:"type:text"`

	// @Description Role of the customer
	// @example customer
	Role string `json:"role" gorm:"not null;default:'customer'"`

	// @Description Whether the customer is active
	// @example false
	IsActive bool `json:"is_active" gorm:"default:true"`

	// @Description Gym owner ID
	// @example 123e4567-e89b-12d3-a456-426614174000
	GymOwnerID uuid.UUID `json:"gym_owner_id" gorm:"type:uuid;not null"`

	// @Description Creation timestamp
	CreatedAt time.Time `json:"created_at"`

	// @Description Last update timestamp
	UpdatedAt time.Time `json:"updated_at"`

	// @Description Soft delete timestamp
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// BeforeSave converts slices to JSON strings before saving to the database
func (c *Customer) BeforeSave(tx *gorm.DB) error {
	if c.HealthConditions != "" {
		c.HealthConditions = string(c.HealthConditions)
	}

	if c.FitnessGoals != "" {
		c.FitnessGoals = string(c.FitnessGoals)
	}

	return nil
}

// AfterFind converts JSON strings to slices after retrieving from the database
func (c *Customer) AfterFind(tx *gorm.DB) error {
	if c.HealthConditions != "" {
		c.HealthConditions = string(c.HealthConditions)
	}

	if c.FitnessGoals != "" {
		c.FitnessGoals = string(c.FitnessGoals)
	}

	return nil
}

// MarshalJSON customizes the JSON response for Customer
func (c Customer) MarshalJSON() ([]byte, error) {
	type Alias Customer
	return json.Marshal(&struct {
		ID string `json:"id"`
		Alias
	}{
		ID:    fmt.Sprintf("cust_%s", c.ID.String()),
		Alias: (Alias)(c),
	})
}

// TableName specifies the table name for the Customer model
func (Customer) TableName() string {
	return "customers"
}

// CreateCustomerRequest represents a request to create a new customer
// @Description Request to create a new customer
type CreateCustomerRequest struct {
	// @Description Full name of the customer
	// @example John Doe
	FullName string `json:"full_name" validate:"required,min=2,max=100"`

	// @Description Email address of the customer
	// @example john.doe@example.com
	Email string `json:"email" validate:"required,email"`

	// @Description Phone number of the customer
	// @example +1234567890
	Phone string `json:"phone" validate:"required,e164"`

	// @Description Gender of the customer
	// @example male
	Gender Gender `json:"gender" validate:"required,oneof=male female other"`

	// @Description Age of the customer
	// @example 25
	Age int `json:"age" validate:"required,min=18,max=100"`

	// @Description Height of the customer in centimeters
	// @example 175
	Height float64 `json:"height" validate:"required,min=100,max=250"`

	// @Description Weight of the customer in kilograms
	// @example 70
	Weight float64 `json:"weight" validate:"required,min=30,max=300"`

	// @Description Health conditions of the customer
	// @example ["Asthma", "Back Pain"]
	HealthConditions string `json:"health_conditions"`

	// @Description Fitness goals of the customer
	// @example ["Weight Loss", "Muscle Gain"]
	FitnessGoals string `json:"fitness_goals"`

	// @Description Profile URL of the customer
	// @example https://example.com/profile.jpg
	ProfileURL string `json:"profile_url" validate:"omitempty,url"`
}

// UpdateCustomerRequest represents a request to update a customer
// @Description Request to update a customer
type UpdateCustomerRequest struct {
	// @Description Full name of the customer
	// @example John Doe
	FullName string `json:"full_name,omitempty" validate:"omitempty,min=2,max=100"`

	// @Description Phone number of the customer
	// @example +1234567890
	Phone string `json:"phone,omitempty" validate:"omitempty,e164"`

	// @Description Email address of the customer
	// @example john.doe@example.com
	Email string `json:"email,omitempty" validate:"omitempty,email"`

	// @Description Gender of the customer
	// @example male
	Gender Gender `json:"gender,omitempty" validate:"omitempty,oneof=male female other unknown"`

	// @Description Age of the customer
	// @example 25
	Age int `json:"age,omitempty" validate:"omitempty,min=18,max=100"`

	// @Description Height of the customer in centimeters
	// @example 175
	Height float64 `json:"height,omitempty" validate:"omitempty,min=100,max=250"`

	// @Description Weight of the customer in kilograms
	// @example 70
	Weight float64 `json:"weight,omitempty" validate:"omitempty,min=30,max=300"`

	// @Description Health conditions of the customer
	// @example ["Asthma", "Back Pain"]
	HealthConditions string `json:"health_conditions,omitempty"`

	// @Description Fitness goals of the customer
	// @example ["Weight Loss", "Muscle Gain"]
	FitnessGoals string `json:"fitness_goals,omitempty"`

	// @Description Profile URL of the customer
	// @example https://example.com/profile.jpg
	ProfileURL string `json:"profile_url,omitempty" validate:"omitempty,url"`
}

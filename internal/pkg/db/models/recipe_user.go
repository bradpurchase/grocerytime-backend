package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// MealUser defines the model for meal_users, which represents a user's usage of a meal
type MealUser struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	MealID uuid.UUID `gorm:"type:uuid;not null"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Meal Meal
	User User
}
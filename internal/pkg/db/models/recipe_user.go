package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// RecipeUser defines the model for recipe_users, which represents a user's usage of a recipe
type RecipeUser struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RecipeID uuid.UUID `gorm:"type:uuid;not null"`
	UserID   uuid.UUID `gorm:"type:uuid;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Recipe Recipe
	User   User
}

package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// RecipeIngredient defines the model for recipe_ingredients
type RecipeIngredient struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RecipeID uuid.UUID `gorm:"type:uuid;not null"`
	Name     string    `gorm:"type:varchar(255);not null;index:idx_recipe_ingredients_name"`
	Amount   int       `gorm:"default:1;not null"`
	Units    *string   `gorm:"type:varchar(20)"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Recipe Recipe
}
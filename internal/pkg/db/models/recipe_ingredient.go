package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// RecipeIngredient defines the model for recipe_ingredients
type RecipeIngredient struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RecipeID uuid.UUID `gorm:"type:uuid;not null;index:idx_recipe_ingredients_recipe_id"`
	Name     string    `gorm:"type:varchar(255);not null"`
	Amount   *float64  `gorm:"default:1"`
	Unit     *string   `gorm:"type:varchar(20)"`
	Notes    *string   `gorm:"type:varchar(255);"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Recipe Recipe
}

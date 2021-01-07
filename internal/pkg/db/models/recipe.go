package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Recipe defines the model for recipes, which represents a global meal object
type Recipe struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index:idx_recipes_user_id"`
	Name     string    `gorm:"type:varchar(255);not null;index:idx_recipes_name"`
	URL      *string   `gorm:"type:varchar(255)"`
	MealType *string   `gorm:"type:varchar(10)"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Ingredients []RecipeIngredient
	Meals       []Meal
}

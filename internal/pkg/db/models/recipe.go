package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Recipe defines the model for recipes
type Recipe struct {
	ID           uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID       uuid.UUID       `gorm:"type:uuid;not null;index:idx_recipes_user_id"`
	Name         string          `gorm:"type:varchar(255);not null;index:idx_recipes_name"`
	Description  *string         `gorm:"type:text"`
	Instructions *datatypes.JSON `gorm:"type:json"`
	MealType     *string         `gorm:"type:varchar(10)"`
	URL          *string         `gorm:"type:varchar(255)"`
	ImageURL     *string         `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Ingredients []RecipeIngredient
	Meals       []Meal
}

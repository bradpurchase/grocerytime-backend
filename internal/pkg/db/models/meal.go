package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Meal defines the model for meals
// A meal is a planned recipe for a user
type Meal struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RecipeID uuid.UUID `gorm:"type:uuid;not null"`
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
	Name     string    `gorm:"type:varchar(255);not null;index:idx_meals_name"`
	MealType *string   `gorm:"type:varchar(10)"`
	Servings int       `gorm:"default:1;not null"`
	Notes    *string   `gorm:"type:text"`
	Date     string    `gorm:"type:varchar(255);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Users  []MealUser
	Recipe Recipe
	User   User
}

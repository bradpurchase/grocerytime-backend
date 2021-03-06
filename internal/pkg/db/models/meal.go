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
	StoreID  uuid.UUID `gorm:"type:uuid;not null"`
	Name     string    `gorm:"type:varchar(255);not null;index:idx_meals_name"`
	MealType *string   `gorm:"type:varchar(10)"`
	Servings int       `gorm:"default:1;not null"`
	Notes    *string   `gorm:"type:text"`
	Date     string    `gorm:"type:varchar(255);not null;index:idx_meals_date"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Users  []MealUser
	Recipe Recipe
	User   User
}

// AfterDelete hook handles deleting associated records after meal is soft-deleted
func (m *Meal) AfterDelete(tx *gorm.DB) (err error) {
	// Delete meal users
	if err := tx.Where("meal_id = ?", m.ID).Delete(&MealUser{}).Error; err != nil {
		return err
	}

	// Remove meal_id assocation on items
	if err := tx.Model(&Item{}).Where("meal_id = ?", m.ID).UpdateColumn("meal_id", nil).Error; err != nil {
		return err
	}

	return
}

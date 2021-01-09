package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type GroceryTripCategory struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	GroceryTripID   uuid.UUID `gorm:"type:uuid;not null;index:idx_grocery_trip_categories_grocery_trip_id"`
	StoreCategoryID uuid.UUID `gorm:"type:uuid;not null;index:idx_grocery_trip_categories_store_category_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	GroceryTrip   GroceryTrip
	StoreCategory StoreCategory
	Item          []Item `gorm:"foreignKey:CategoryID"`
}

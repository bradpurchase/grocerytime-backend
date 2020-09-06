package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type GroceryTripCategory struct {
	ID                uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	GroceryTripID     uuid.UUID `gorm:"type:uuid;not null"`
	StoreCategoryID   uuid.UUID `gorm:"type:uuid;not null"`
	StoreCategoryName string    // Supports join column (see: internal/pkg/gql/types/grocery_trip.go:41)

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Associations
	GroceryTrip   GroceryTrip
	StoreCategory StoreCategory
	Item          []Item
}

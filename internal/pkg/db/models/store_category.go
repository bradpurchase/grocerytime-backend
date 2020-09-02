package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type StoreCategory struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	StoreID    uuid.UUID `gorm:"type:uuid;not null"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Associations
	Store    Store
	Category Category
	Items    []Item
}

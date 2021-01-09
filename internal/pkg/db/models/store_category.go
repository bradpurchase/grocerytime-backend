package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type StoreCategory struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StoreID uuid.UUID `gorm:"type:uuid;not null;index:idx_store_categories_store_id"`
	Name    string    `gorm:"type:varchar(100);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Store Store
}

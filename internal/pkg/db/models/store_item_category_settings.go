package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type StoreItemCategorySettings struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StoreID uuid.UUID `gorm:"type:uuid;not null;index:idx_store_category_items_store_id"`
	Items   datatypes.JSON

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

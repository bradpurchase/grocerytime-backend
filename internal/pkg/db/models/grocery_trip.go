package models

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type GroceryTrip struct {
	ID                 uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StoreID            uuid.UUID `gorm:"type:uuid;not null;index:idx_store_categories_store_id"`
	Name               string    `gorm:"type:varchar(100);not null"`
	Completed          bool      `gorm:"default:false;not null"`
	CopyRemainingItems bool      `gorm:"default:false;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Store Store
	Items []Item
}

// BeforeCreate hook is triggered before a trip is created
func (g *GroceryTrip) BeforeCreate(tx *gorm.DB) (err error) {
	// If this store has already had a trip with this name, affix a count to it to make it unique
	var count int64
	name := fmt.Sprintf("%%%s%%", g.Name) // LIKE '%Oct 08, 2020%'
	if err := tx.Model(&GroceryTrip{}).Where("name LIKE ? AND store_id = ?", name, g.StoreID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		g.Name = fmt.Sprintf("%s (%d)", g.Name, count+1)
	}
	return
}

package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Store struct {
	ID     uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ListID uuid.UUID `gorm:"type:uuid;not null"`
	Name   string    `gorm:"type:varchar(100);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Associations
	List List
}

// AfterCreate hook to automatically create default store_categories
func (s *Store) AfterCreate(tx *gorm.DB) (err error) {
	// Fetch the global store of categories
	categories := []Category{}
	if err := tx.Find(&categories).Error; err != nil {
		return err
	}

	for i := range categories {
		storeCategory := &StoreCategory{StoreID: s.ID, CategoryID: categories[i].ID}
		if err := tx.Create(&storeCategory).Error; err != nil {
			return err
		}
	}

	return nil
}

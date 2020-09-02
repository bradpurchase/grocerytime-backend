package models

import (
	"time"

	"github.com/jinzhu/gorm"

	uuid "github.com/satori/go.uuid"
)

type Store struct {
	ID     uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();not null"`
	Name   string    `gorm:"type:varchar(100);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Associations
	StoreUsers   []StoreUser
	Item         []Item
	GroceryTrips []GroceryTrip
}

// AfterCreate hook to automatically create some associated records
func (s *Store) AfterCreate(tx *gorm.DB) (err error) {
	// Create default store user (creator)
	creator := true
	active := true
	storeUser := StoreUser{
		StoreID: s.ID,
		UserID:  s.UserID,
		Creator: &creator,
		Active:  &active,
	}
	if err := tx.Create(&storeUser).Error; err != nil {
		return err
	}

	// Create categories for the default store
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

	// Create default grocery trip
	trip := GroceryTrip{
		StoreID:            s.ID,
		Name:               "Trip 1",
		Completed:          false,
		CopyRemainingItems: false,
	}
	if err := tx.Create(&trip).Error; err != nil {
		return err
	}

	return nil
}

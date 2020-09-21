package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Store struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();not null"`
	Name   string    `gorm:"type:varchar(100);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	StoreUsers   []StoreUser
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

	categories := fetchCategories()
	for i := range categories {
		storeCategory := &StoreCategory{StoreID: s.ID, Name: categories[i]}
		if err := tx.Create(&storeCategory).Error; err != nil {
			return err
		}
	}

	// Create default grocery trip
	currentTime := time.Now()
	tripName := currentTime.Format("Jan 02, 2006")
	trip := GroceryTrip{
		StoreID:            s.ID,
		Name:               tripName,
		Completed:          false,
		CopyRemainingItems: false,
	}
	if err := tx.Create(&trip).Error; err != nil {
		return err
	}

	return nil
}

func fetchCategories() [20]string {
	categories := [20]string{
		"Produce",
		"Bakery",
		"Meat",
		"Seafood",
		"Dairy",
		"Cereal",
		"Baking",
		"Dry Goods",
		"Canned Goods",
		"Frozen Foods",
		"Cleaning",
		"Paper Products",
		"Beverages",
		"Candy & Snacks",
		"Condiments",
		"Personal Care",
		"Baby",
		"Alcohol",
		"Pharmacy",
		"Misc.",
	}
	return categories
}

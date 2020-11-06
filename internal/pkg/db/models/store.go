package models

import (
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
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
	tripName := currentTime.Format("Jan 2, 2006")
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

// AfterDelete hook handles deleting associated records after store is deleted
func (s *Store) AfterDelete(tx *gorm.DB) (err error) {
	// Delete items associated with this store
	var trips []GroceryTrip
	if err := tx.Where("store_id = ?", s.ID).Find(&trips).Error; err != nil {
		return err
	}
	for i := range trips {
		if err := tx.Where("grocery_trip_id = ?", trips[i].ID).Delete(&Item{}).Error; err != nil {
			return err
		}
	}

	// Delete GroceryTrip records associated with this store
	if err := tx.Where("store_id = ?", s.ID).Delete(&GroceryTrip{}).Error; err != nil {
		return err
	}

	// Delete StoreUser records associated with this store
	var storeUsers []StoreUser
	if err := tx.Where("store_id = ? AND active = ?", s.ID, true).Find(&storeUsers).Error; err != nil {
		return err
	}

	// Send notification to store users about this store being deleted
	var emails []string
	for i := range storeUsers {
		var user User
		userQuery := tx.
			Select("email").
			Where("id = ?", storeUsers[i].UserID).
			Find(&user).
			Error
		if err := userQuery; err != nil {
			return err
		}
		emails = append(emails, user.Email)
	}

	_, e := mailer.SendStoreDeletedEmail(s.Name, emails)
	if e != nil {
		return e
	}

	if err := tx.Where("store_id = ?", s.ID).Delete(&StoreUser{}).Error; err != nil {
		return err
	}

	return
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

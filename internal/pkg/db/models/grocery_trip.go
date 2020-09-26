package models

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type GroceryTrip struct {
	ID                 uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StoreID            uuid.UUID `gorm:"type:uuid;not null"`
	Name               string    `gorm:"type:varchar(100);not null"`
	Completed          bool      `gorm:"default:false;not null"`
	CopyRemainingItems bool      `gorm:"default:false;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Store Store
}

// AfterUpdate hook is triggered after a trip is updated, such as in trips.UpdateTrip
func (g *GroceryTrip) AfterUpdate(tx *gorm.DB) (err error) {
	if g.Completed {
		var tripsCount int64
		if err := tx.Model(&GroceryTrip{}).Where("store_id = ?", g.StoreID).Count(&tripsCount).Error; err != nil {
			return err
		}

		// Create the next trip for the user
		currentTime := time.Now()
		newTripName := currentTime.Format("Jan 02, 2006")
		newTrip := &GroceryTrip{StoreID: g.StoreID, Name: newTripName}
		if err := tx.Create(&newTrip).Error; err != nil {
			return err
		}

		// If the completed trip was configured to copy its remaining items
		// over to the next trip, perform this operation - otherwise, mark
		// each item in the completed trip as completed
		completed := true
		columns := Item{Completed: &completed}
		if g.CopyRemainingItems {
			// Duplicate the category associated with each item
			var remainingItems []Item
			if err := tx.Where("grocery_trip_id = ? AND completed = ?", g.ID, false).Find(&remainingItems).Error; err != nil {
				return err
			}
			for i := range remainingItems {
				// Retrieve the store category associated with the previous trip and use it
				// to create a duplicate grocery trip category in new trip
				// (note: use FindOrCreate to handle multiple items in same category)
				storeCategory := &StoreCategory{}
				storeCategoryQuery := tx.
					Select("store_categories.id, store_categories.name").
					Joins("INNER JOIN grocery_trip_categories ON grocery_trip_categories.store_category_id = store_categories.id").
					Where("grocery_trip_categories.id = ?", remainingItems[i].CategoryID).
					Find(&storeCategory).
					Error
				if err := storeCategoryQuery; err != nil {
					return err
				}
				// Note: uses FirstOrCreate to handle the case where there are multiple items in same category
				// that need to be moved over to the next trip
				groceryTripCategory := GroceryTripCategory{GroceryTripID: newTrip.ID, StoreCategoryID: storeCategory.ID}
				if err := tx.Where(groceryTripCategory).FirstOrCreate(&groceryTripCategory).Error; err != nil {
					return err
				}
				columns = Item{GroceryTripID: newTrip.ID, CategoryID: &groceryTripCategory.ID}
			}
		}
		updateItemsQuery := tx.
			Model(&Item{}).
			Where("grocery_trip_id = ? AND completed = ?", g.ID, false).
			UpdateColumns(columns).
			Error
		if err := updateItemsQuery; err != nil {
			return err
		}
	}
	return nil
}

// BeforeCreate hook is triggered before a trip is created
func (g *GroceryTrip) BeforeCreate(tx *gorm.DB) (err error) {
	// If this store has already had a trip with this name, affix a count to it to make it unique
	var count int64
	if err := tx.Model(&GroceryTrip{}).Where("name = ? AND store_id = ?", g.Name, g.StoreID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		g.Name = fmt.Sprintf("%s (%d)", g.Name, count + 1)
	}
	return
}

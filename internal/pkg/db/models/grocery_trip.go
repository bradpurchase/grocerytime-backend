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
	Items []Item
}

// AfterUpdate hook is triggered after a trip is updated, such as in trips.UpdateTrip
func (g *GroceryTrip) AfterUpdate(tx *gorm.DB) (err error) {
	if g.Completed {
		// Create the next trip for the user
		var tripsCount int64
		if err := tx.Model(&GroceryTrip{}).Where("store_id = ?", g.StoreID).Count(&tripsCount).Error; err != nil {
			return err
		}
		currentTime := time.Now()
		newTripName := currentTime.Format("Jan 02, 2006")
		newTrip := &GroceryTrip{StoreID: g.StoreID, Name: newTripName}
		if err := tx.Create(&newTrip).Error; err != nil {
			return err
		}

		// Handle copy operation the completed trip was configured to copy its
		// remaining items over to the next trip

		if g.CopyRemainingItems {
			// Duplicate the category associated with each item
			var remainingItems []Item
			if err := tx.Where("grocery_trip_id = ? AND completed = ?", g.ID, false).Find(&remainingItems).Error; err != nil {
				return err
			}

			var newItems []Item
			for i := range remainingItems {
				// Retrieve the store category associated with the previous trip and use it
				// to create a duplicate grocery trip category in new trip
				// (note: use FindOrCreate to avoid dupe categories)
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

				// Copy old item to new item and update values
				newItem := remainingItems[i]
				newItem.ID = uuid.Nil
				newItem.GroceryTripID = newTrip.ID
				newItem.CategoryID = &groceryTripCategory.ID
				newItem.CreatedAt = time.Now()
				newItem.UpdatedAt = time.Now()
				newItems = append(newItems, newItem)
			}

			// Batch insert items in new trip
			if err := tx.Create(&newItems).Error; err != nil {
				return err
			}
		}

		// Mark each item in the old trip as completed
		//
		// This uses UpdateColumn to avoid hooks
		// (https://gorm.io/docs/update.html#Without-Hooks-Time-Tracking)
		updateItemsQuery := tx.
			Model(&Item{}).
			Where("grocery_trip_id = ? AND completed = ?", g.ID, false).
			UpdateColumn("completed", true).
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
	name := fmt.Sprintf("%%%s%%", g.Name) // LIKE '%Oct 08, 2020%'
	if err := tx.Model(&GroceryTrip{}).Where("name LIKE ? AND store_id = ?", name, g.StoreID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		g.Name = fmt.Sprintf("%s (%d)", g.Name, count+1)
	}
	return
}

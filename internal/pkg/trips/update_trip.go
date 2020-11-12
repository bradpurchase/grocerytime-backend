package trips

import (
	"errors"
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// UpdateTrip updates a grocery trip with the given args by tripID
func UpdateTrip(args map[string]interface{}) (interface{}, error) {
	trip := models.GroceryTrip{}
	if err := db.Manager.Where("id = ?", args["tripId"]).First(&trip).Error; err != nil {
		return nil, errors.New("trip does not exist")
	}
	if args["name"] != nil {
		trip.Name = args["name"].(string)
	}
	if args["completed"] != nil {
		trip.Completed = args["completed"].(bool)
	}
	if args["copyRemainingItems"] != nil {
		trip.CopyRemainingItems = args["copyRemainingItems"].(bool)
	}
	if err := db.Manager.Save(&trip).Error; err != nil {
		return nil, err
	}

	// If the trip was completed, create the next trip for the user
	if trip.Completed {
		db.Manager.Transaction(func(tx *gorm.DB) error {
			var tripsCount int64
			tripsCountQuery := tx.
				Model(&models.GroceryTrip{}).
				Where("store_id = ?", trip.StoreID).
				Count(&tripsCount).
				Error
			if err := tripsCountQuery; err != nil {
				return err
			}

			var newTripName string
			// If a newTripName argument is passed, use it instead of creating one from the
			// current date in server time; typically this argument will be passed from
			// the app as the current date from the device
			if args["newTripName"] != nil {
				newTripName = args["newTripName"].(string)
			} else {
				currentTime := time.Now()
				newTripName = currentTime.Format("Jan 2, 2006")
			}
			newTrip := &models.GroceryTrip{StoreID: trip.StoreID, Name: newTripName}
			if err := tx.Create(&newTrip).Error; err != nil {
				return err
			}

			if trip.CopyRemainingItems {
				// Duplicate the category associated with each item
				var remainingItems []models.Item
				if err := tx.Where("grocery_trip_id = ? AND completed = ?", trip.ID, false).Find(&remainingItems).Error; err != nil {
					return err
				}

				var newItems []models.Item
				for i := range remainingItems {
					// Retrieve the store category associated with the previous trip and use it
					// to create a duplicate grocery trip category in new trip
					// (note: use FindOrCreate to avoid dupe categories)
					storeCategory := &models.StoreCategory{}
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
					groceryTripCategory := models.GroceryTripCategory{
						GroceryTripID:   newTrip.ID,
						StoreCategoryID: storeCategory.ID,
					}
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
				Model(&models.Item{}).
				Where("grocery_trip_id = ? AND completed = ?", trip.ID, false).
				UpdateColumn("completed", true).
				Error
			if err := updateItemsQuery; err != nil {
				return err
			}

			return nil
		})
	}

	return trip, nil
}

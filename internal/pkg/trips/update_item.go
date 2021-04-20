package trips

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// UpdateItem updates an item by itemID
func UpdateItem(args map[string]interface{}) (interface{}, error) {
	item := &models.Item{}
	if err := db.Manager.Where("id = ?", args["itemId"]).First(&item).Error; err != nil {
		return nil, err
	}

	if args["name"] != nil {
		item.Name = args["name"].(string)
	}
	if args["completed"] != nil {
		completed := args["completed"].(bool)
		item.Completed = &completed
		//item.Position = GetNewPosition(db, item.GroceryTripID, completed)
	}
	if args["quantity"] != nil {
		item.Quantity = args["quantity"].(int)
	}
	if args["position"] != nil {
		item.Position = args["position"].(int)
	}
	if args["notes"] != nil {
		notes := args["notes"].(string)
		item.Notes = &notes
	}
	if args["storeCategoryId"] != nil {
		storeCategoryID, _ := uuid.FromString(args["storeCategoryId"].(string))
		groceryTripCategory := models.GroceryTripCategory{
			GroceryTripID:   item.GroceryTripID,
			StoreCategoryID: storeCategoryID,
		}
		// This returns nothing unless the category is already present in the trip,
		// so we do a FirstOrCreate here
		if err := db.Manager.Where(groceryTripCategory).FirstOrCreate(&groceryTripCategory).Error; err != nil {
			return nil, err
		}
		item.CategoryID = &groceryTripCategory.ID
	}
	if err := db.Manager.Save(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// GetNewPosition gets the new position of an updated item
func GetNewPosition(tripID uuid.UUID, completed bool) int {
	newPosition := 1
	if completed {
		// If the item was marked completed, move to the bottom of the store
		// The BeforeUpdate hook on items will handle reordering the items around it
		bottomItem := &models.Item{}
		db.Manager.
			Select("position").
			Where("grocery_trip_id = ?", tripID).
			Order("position DESC").
			Limit(1).
			Find(&bottomItem)
		newPosition = bottomItem.Position
	}
	return newPosition
}

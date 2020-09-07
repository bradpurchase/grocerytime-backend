package trips

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// UpdateItem updates an item by itemID
func UpdateItem(db *gorm.DB, args map[string]interface{}) (interface{}, error) {
	item := &models.Item{}
	if err := db.Where("id = ?", args["itemId"]).First(&item).Error; err != nil {
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
	if args["categoryId"] != nil {
		categoryID := args["categoryId"].(uuid.UUID)
		item.CategoryID = &categoryID
	}
	if err := db.Save(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// GetNewPosition gets the new position of an updated item
func GetNewPosition(db *gorm.DB, tripID uuid.UUID, completed bool) int {
	newPosition := 1
	if completed {
		// If the item was marked completed, move to the bottom of the store
		// The BeforeUpdate hook on items will handle reordering the items around it
		bottomItem := &models.Item{}
		db.
			Select("position").
			Where("grocery_trip_id = ?", tripID).
			Order("position DESC").
			Limit(1).
			Find(&bottomItem)
		newPosition = bottomItem.Position
	}
	return newPosition
}

package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
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
		item.Completed = completed
		item.Position = GetNewPosition(db, item.ListID, completed)
	}
	if args["quantity"] != nil {
		item.Quantity = args["quantity"].(int)
	}
	if args["position"] != nil {
		item.Position = args["position"].(int)
	}
	if err := db.Save(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// GetNewPosition gets the new position of an updated item
func GetNewPosition(db *gorm.DB, listID uuid.UUID, completed bool) int {
	// Reorder to the bottom of the list
	newPosition, err := DetermineListPosition("top", db, listID)
	if completed {
		newPosition, err = DetermineListPosition("bottom", db, listID)
	}
	if err != nil {
		return 0
	}
	return newPosition
}

package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
)

// UpdateItem updates an item by itemID
func UpdateItem(db *gorm.DB, args map[string]interface{}) (interface{}, error) {
	item := &models.Item{}
	if err := db.Where("id = ?", args["itemId"]).First(&item).Error; err != nil {
		return nil, err
	}

	if args["completed"] != nil {
		item.Completed = args["completed"].(bool)
	}
	if args["quantity"] != nil {
		item.Quantity = args["quantity"].(int)
	}
	if err := db.Save(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

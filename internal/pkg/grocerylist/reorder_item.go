package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
)

// ReorderItem handles the reordering of an item in the list by taking the
// item ID and the new position. It returns he reordered list.
func ReorderItem(db *gorm.DB, itemID interface{}, position int) (*models.GroceryTrip, error) {
	trip := &models.GroceryTrip{}
	item := &models.Item{}
	if err := db.Where("id = ?", itemID).First(&item).Error; err != nil {
		return trip, err
	}
	item.Position = position
	if err := db.Save(&item).Error; err != nil {
		return trip, err
	}
	if err := db.Where("id = ?", item.GroceryTripID).Find(&trip).Error; err != nil {
		return trip, err
	}
	return trip, nil
}

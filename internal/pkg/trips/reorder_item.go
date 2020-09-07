package trips

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"gorm.io/gorm"
)

// ReorderItem handles the reordering of an item by taking the
// item ID and the new position. It returns the reordered trip object.
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

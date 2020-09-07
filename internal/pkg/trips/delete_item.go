package trips

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
)

// DeleteItem deletes an item from a trip and handles trip category cleanup
// (i.e. if this is the last item in a trip category, it deletes the trip category)
func DeleteItem(db *gorm.DB, itemID interface{}) (interface{}, error) {
	item := models.Item{}
	if err := db.Where("id = ?", itemID).First(&item).Error; err != nil {
		return nil, errors.New("item not found")
	}

	categoryID := item.CategoryID
	if err := db.Delete(&item).Error; err != nil {
		return nil, err
	}

	// If this was the last item in this trip category, delete the trip category too
	remainingItemsCount := 0
	if err := db.Model(&models.Item{}).Where("category_id = ?", categoryID).Count(&remainingItemsCount).Error; err != nil {
		return nil, err
	}
	if remainingItemsCount == 0 {
		if err := db.Where("id = ?", categoryID).Delete(&models.GroceryTripCategory{}).Error; err != nil {
			return nil, err
		}
	}

	return item, nil
}

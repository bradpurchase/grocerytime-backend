package trips

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// DeleteItem deletes an item from a trip and handles trip category cleanup
// (i.e. if this is the last item in a trip category, it deletes the trip category)
func DeleteItem(itemID interface{}) (deletedItem models.Item, err error) {
	item := models.Item{}
	if err := db.Manager.Where("id = ?", itemID).First(&item).Error; err != nil {
		return deletedItem, errors.New("item not found")
	}

	categoryID := item.CategoryID
	if err := db.Manager.Delete(&item).Error; err != nil {
		return deletedItem, err
	}

	// If this was the last item in this trip category, delete the trip category too
	var remainingItemsCount int64
	if err := db.Manager.Model(&models.Item{}).Where("category_id = ?", categoryID).Count(&remainingItemsCount).Error; err != nil {
		return deletedItem, err
	}
	if remainingItemsCount == 0 {
		if err := db.Manager.Where("id = ?", categoryID).Delete(&models.GroceryTripCategory{}).Error; err != nil {
			return deletedItem, err
		}
	}

	return item, nil
}

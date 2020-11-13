package trips

import (
	"errors"
	"time"

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

	// Touch the GroceryTrip record to update its updated_at timestamp
	updateTripQuery := db.Manager.
		Model(&models.GroceryTrip{}).
		Where("id = ?", item.GroceryTripID).
		Update("updated_at", time.Now()).
		Error
	if err := updateTripQuery; err != nil {
		return deletedItem, err
	}

	// If this was the last item in this trip category, delete the trip category too
	var remainingItemsCount int64
	categoryItemsCountQuery := db.Manager.
		Model(&models.Item{}).
		Where("category_id = ?", categoryID).
		Count(&remainingItemsCount).
		Error
	if err := categoryItemsCountQuery; err != nil {
		return deletedItem, err
	}
	if remainingItemsCount == 0 {
		deleteCategoryQuery := db.Manager.
			Where("id = ?", categoryID).
			Delete(&models.GroceryTripCategory{}).
			Error
		if err := deleteCategoryQuery; err != nil {
			return deletedItem, err
		}
	}

	return item, nil
}

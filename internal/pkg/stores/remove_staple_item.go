package stores

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// RemoveStapleItem saves an item as a staple in the store ID provided
func RemoveStapleItem(itemID uuid.UUID) (staple models.StoreStapleItem, err error) {
	var item models.Item
	if err := db.Manager.Select("staple_item_id").Where("id = ?", itemID).First(&item).Error; err != nil {
		return staple, err
	}

	if err := db.Manager.Where("id = ?", item.StapleItemID).Delete(&staple).Error; err != nil {
		return staple, err
	}

	// Update the item to dissociate staple_item_id
	if err := db.Manager.Model(&item).Where("staple_item_id = ?", item.StapleItemID).Update("staple_item_id", nil).Error; err != nil {
		return staple, err
	}

	return staple, nil
}

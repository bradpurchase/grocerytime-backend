package stores

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// RemoveStapleItem dissocates any staple_item_id from items for this staple item and deletes the staple item
func RemoveStapleItem(itemID uuid.UUID) (staple models.StoreStapleItem, err error) {
	var item models.Item
	if err := db.Manager.Select("staple_item_id").Where("id = ?", itemID).First(&item).Error; err != nil {
		return staple, err
	}
	if err := db.Manager.Where("id = ?", item.StapleItemID).Delete(&staple).Error; err != nil {
		return staple, err
	}

	updateQuery := db.Manager.
		Model(&item).
		Where("staple_item_id = ?", item.StapleItemID).
		UpdateColumn("staple_item_id", nil).
		Error
	if err := updateQuery; err != nil {
		return staple, err
	}

	return staple, nil
}

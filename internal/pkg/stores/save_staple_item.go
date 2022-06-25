package stores

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// SaveStapleItem saves an item as a staple in the store ID provided
func SaveStapleItem(storeID uuid.UUID, itemID uuid.UUID) (staple models.StoreStapleItem, err error) {
	var item models.Item
	if err := db.Manager.Where("id = ?", itemID).First(&item).Error; err != nil {
		return staple, err
	}

	stapleItem := models.StoreStapleItem{StoreID: storeID, Name: item.Name}
	if err := db.Manager.Where(stapleItem).FirstOrCreate(&stapleItem).Error; err != nil {
		return staple, err
	}

	// This makes for quick lookup and setting/unsetting items as staples
	if err := db.Manager.Model(&item).UpdateColumn("staple_item_id", &stapleItem.ID).Error; err != nil {
		return stapleItem, err
	}

	return stapleItem, nil
}

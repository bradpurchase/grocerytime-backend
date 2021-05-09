package stores

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// RemoveStapleItem saves an item as a staple in the store ID provided
func RemoveStapleItem(stapleItemID uuid.UUID) (staple models.StoreStapleItem, err error) {
	if err := db.Manager.Where("id = ?", stapleItemID).First(&staple).Error; err != nil {
		return staple, err
	}

	// This makes for quick lookup and setting/unsetting items as staples
	if err := db.Manager.Delete(&staple).Error; err != nil {
		return staple, err
	}
	return staple, nil
}

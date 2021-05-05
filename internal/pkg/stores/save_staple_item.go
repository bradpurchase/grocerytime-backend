package stores

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// SaveStapleItem saves an item as a staple in the store ID provided
func SaveStapleItem(userID uuid.UUID, storeID uuid.UUID, name string) (staple models.StoreStapleItem, err error) {
	userExistsQuery := db.Manager.
		First(&models.StoreUser{}).
		Where("user_id = ? AND store_id = ?", userID, storeID).
		Error
	if errors.Is(userExistsQuery, gorm.ErrRecordNotFound) {
		return staple, errors.New("user does not belong to this store")
	}

	stapleItem := models.StoreStapleItem{StoreID: storeID, Name: name}
	if err := db.Manager.Where(stapleItem).FirstOrCreate(&stapleItem).Error; err != nil {
		return staple, err
	}
	return stapleItem, nil
}

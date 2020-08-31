package stores

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// RetrieveStoreForList retrieves the store for the given listID
func RetrieveStoreForList(db *gorm.DB, listID uuid.UUID) (interface{}, error) {
	store := &models.Store{}
	if err := db.Where("list_id = ?", listID).First(&store).Error; err != nil {
		return nil, err
	}
	return store, nil
}

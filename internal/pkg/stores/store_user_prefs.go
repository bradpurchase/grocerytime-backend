package stores

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// RetrieveStoreUserPrefs updates store user preferences
func RetrieveStoreUserPrefs(storeUserID uuid.UUID) (sup models.StoreUserPreference, err error) {
	query := db.Manager.
		Where("store_user_id = ?", storeUserID).
		First(&sup).
		Error
	if err := query; err != nil {
		return sup, err
	}
	return sup, nil
}

package stores

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// UpdateStoreUserPrefs updates store user preferences
func UpdateStoreUserPrefs(storeUserID uuid.UUID, args map[string]interface{}) (sup models.StoreUserPreference, err error) {
	query := db.Manager.
		Where("store_user_id = ?", storeUserID).
		First(&sup).
		Error
	if err := query; err != nil {
		return sup, err
	}

	if args["defaultStore"] != nil {
		sup.DefaultStore = args["defaultStore"].(bool)
	}
	if args["notifications"] != nil {
		sup.Notifications = args["notifications"].(bool)
	}

	if err := db.Manager.Save(&sup).Error; err != nil {
		return sup, err
	}
	return sup, nil
}

package stores

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// UpdateStoreForUser updates a store for the given userID with the provided args
func UpdateStoreForUser(db *gorm.DB, userID uuid.UUID, args map[string]interface{}) (interface{}, error) {
	store := &models.Store{}
	if err := db.Where("id = ? AND user_id = ?", args["storeId"], userID).First(&store).Error; err != nil {
		return nil, err
	}

	oldName := store.Name
	if args["name"] != nil {
		store.Name = args["name"].(string)
	}
	if err := db.Save(&store).Error; err != nil {
		return nil, err
	}

	// Finally, send an email to the users of this store about this update (excluding the creator)
	if oldName != args["name"] {
		rows, err := db.Raw("SELECT u.email FROM store_users AS su INNER JOIN users AS u ON su.user_id = u.id WHERE su.store_id = ? AND su.creator = ? ORDER BY su.created_at DESC", store.ID, false).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var email string
		for rows.Next() {
			rows.Scan(&email)
			_, mailErr := mailer.SendStoreRenamedEmail(oldName, store.Name, email)
			if mailErr != nil {
				return nil, mailErr
			}
		}
	}

	return store, nil
}

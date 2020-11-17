package models

import (
	"time"

	"gorm.io/gorm"

	uuid "github.com/satori/go.uuid"
)

type StoreUserPreference struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StoreUserID   uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	DefaultStore  bool      `gorm:"default:false;not null"`
	Notifications bool      `gorm:"default:true;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// AfterUpdate hook handles some cleanup operations after updating store user prefs
func (sup *StoreUserPreference) AfterUpdate(tx *gorm.DB) (err error) {
	if sup.DefaultStore {
		// Find all other stores belonging to this user and mark them as default_store false
		var storeUsers []StoreUser
		storeUsersQuery := tx.
			Select("id").
			Where("user_id IN (?)", tx.Model(&StoreUser{}).Select("user_id").Where("id = ?", sup.StoreUserID)).
			Where("id NOT IN (?)", sup.StoreUserID).
			Find(&storeUsers).
			Error
		if err := storeUsersQuery; err != nil {
			return err
		}

		// Collect store user IDs and use them to update all other records to remove default flag
		var storeUserIds []uuid.UUID
		for i := range storeUsers {
			storeUserIds = append(storeUserIds, storeUsers[i].ID)
		}
		if len(storeUserIds) > 0 {
			updateQuery := tx.
				Model(&StoreUserPreference{}).
				Where("store_user_id IN (?)", storeUserIds).
				Update("default_store", false).
				Error
			if err := updateQuery; err != nil {
				return err
			}
		}
	}
	return nil
}

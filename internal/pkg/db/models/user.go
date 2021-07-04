package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                       uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email                    string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password                 string     `gorm:"not null"`
	Name                     string     `gorm:"type:varchar(100)"`
	PasswordResetToken       *uuid.UUID `gorm:"type:uuid"`
	PasswordResetTokenExpiry *time.Time
	SiwaID                   *string `gorm:"type:varchar(255);uniqueIndex"`

	LastSeenAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Associations
	Stores []Store
	Tokens []AuthToken
}

// BeforeDelete handles removing associated data before a user account is deleted
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	// Hard-delete auth tokens
	var authTokens []AuthToken
	if err := tx.Unscoped().Where("user_id = ?", u.ID).Delete(&authTokens).Error; err != nil {
		return err
	}

	// Hard-delete devices
	var devices []Device
	if err := tx.Unscoped().Where("user_id = ?", u.ID).Delete(&devices).Error; err != nil {
		return err
	}

	// Delete store users
	var storeUsers []StoreUser
	if err := tx.Unscoped().Where("user_id = ?", u.ID).Delete(&storeUsers).Error; err != nil {
		return err
	}

	// Delete stores
	// The Store model has an AfterDelete hook which handles deleting associated
	// records after the store is deleted
	var userStores []Store
	if err := tx.Unscoped().Where("user_id = ?", u.ID).Delete(&userStores).Error; err != nil {
		return err
	}

	// TODO: Delete meals/recipes
	return nil
}

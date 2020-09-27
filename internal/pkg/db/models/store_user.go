package models

import (
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	"gorm.io/gorm"

	uuid "github.com/satori/go.uuid"
)

type StoreUser struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StoreID      uuid.UUID `gorm:"type:uuid;not null"`
	UserID       uuid.UUID `gorm:"type:uuid"`
	Email        string    `gorm:"type:varchar(100)"`
	Creator      *bool     `gorm:"default:false;not null"`
	Active       *bool     `gorm:"default:true;not null"`
	DefaultStore bool      `gorm:"default:false;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	//Store Store
	//User  User
}

// BeforeCreate hook to handle setting some properties before adding the record
func (su *StoreUser) BeforeCreate(tx *gorm.DB) (err error) {
	// Determine whether or not this store should be considered the default
	// store for the user. This should be the case if it's the first store they're a part of
	var numStores int64
	if err := tx.Model(&StoreUser{}).Where("user_id = ? AND active = ?", su.UserID, true).Count(&numStores).Error; err != nil {
		return err
	}
	su.DefaultStore = numStores == 0
	return
}

// BeforeUpdate hook to handle setting DefaultStore if user was invited and has now joined their first store
func (su *StoreUser) BeforeUpdate(tx *gorm.DB) (err error) {
	var numStores int64
	if err := tx.Model(&StoreUser{}).Where("user_id = ? AND active = ?", su.UserID, true).Count(&numStores).Error; err != nil {
		return err
	}
	su.DefaultStore = numStores == 1
	return
}

// AfterCreate hook to handle sending an invite email to a new StoreUser if the
// email column is not empty (i.e. store invitation by another user)
func (su *StoreUser) AfterCreate(tx *gorm.DB) (err error) {
	if len(su.Email) > 0 {
		store := Store{}
		if err := tx.Select("name, user_id").Where("id = ?", su.StoreID).First(&store).Error; err != nil {
			return err
		}

		// Retrieve the name of the user who created this list by using store.UserID
		creatorUser := User{}
		if err := tx.Select("name").Where("id = ?", store.UserID).First(&creatorUser).Error; err != nil {
			return err
		}
		_, err := mailer.SendStoreInvitationEmail(store.Name, su.Email, creatorUser.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

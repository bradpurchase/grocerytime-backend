package models

import (
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	"gorm.io/gorm"

	uuid "github.com/satori/go.uuid"
)

type StoreUser struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StoreID uuid.UUID `gorm:"type:uuid;not null;index:idx_store_users_store_id"`
	UserID  uuid.UUID `gorm:"type:uuid"`
	Email   string    `gorm:"type:varchar(100)"`
	Creator *bool     `gorm:"default:false;not null"`
	Active  *bool     `gorm:"default:true;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Preferences StoreUserPreference
	Store       Store
	User        User
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
	} else {
		// handle creating a StoreUserPreference record for new StoreUser records
		// that don't have an email attached (invite case)
		prefs := StoreUserPreference{StoreUserID: su.ID}
		if err := tx.Create(&prefs).Error; err != nil {
			return err
		}
	}
	return
}

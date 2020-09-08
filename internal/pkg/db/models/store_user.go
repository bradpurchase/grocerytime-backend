package models

import (
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	"gorm.io/gorm"

	uuid "github.com/satori/go.uuid"
)

type StoreUser struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StoreID uuid.UUID `gorm:"type:uuid;not null"`
	UserID  uuid.UUID `gorm:"type:uuid"`
	Email   string    `gorm:"type:varchar(100)"`
	Creator *bool     `gorm:"default:false;not null"`
	Active  *bool     `gorm:"default:true;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	Store Store
	User  User
}

// AfterCreate hook to handle sending an invite email to a new StoreUser if the
// email column is not empty (i.e. store invitation by another user)
func (lu *StoreUser) AfterCreate(tx *gorm.DB) (err error) {
	if len(lu.Email) > 0 {
		store := Store{}
		if err := tx.Select("name").Where("id = ?", lu.StoreID).First(&store).Error; err != nil {
			return err
		}
		_, err := mailer.SendStoreInvitationEmail(store.Name, lu.Email)
		if err != nil {
			return err
		}
	}
	return nil
}

package models

import (
	"time"

	// Postgres dialect for GORM
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

type ListUser struct {
	ID      uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ListID  uuid.UUID `gorm:"type:uuid;not null"`
	UserID  uuid.UUID `gorm:"type:uuid"`
	Email   string    `gorm:"type:varchar(100)"`
	Creator bool      `gorm:"default:false;not null"`
	Active  bool      `gorm:"default:true;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations
	List      List
	User      User
	ListUsers []ListUser
}

// AfterCreate hook to handle sending an invite email to a new ListUser if the
// email column is not empty (i.e. list invitation by another user)
func (lu *ListUser) AfterCreate(tx *gorm.DB) (err error) {
	if len(lu.Email) > 0 {
		//TODO lu.List.Name doesnt work :( prolly need to fetch the list separately
		_, err := mailer.SendListInvitationEmail(lu.List.Name, lu.Email)
		if err != nil {
			return err
		}
	}
	return nil
}

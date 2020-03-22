package models

import (
	"time"

	"github.com/jinzhu/gorm"

	uuid "github.com/satori/go.uuid"
)

type List struct {
	ID     uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();not null"`
	Name   string    `gorm:"type:varchar(100);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations
	ListUsers    []ListUser
	GroceryTrips []GroceryTrip
}

// AfterCreate hook to automatically create a ListUser record
// for the user who is creating the list
func (l *List) AfterCreate(tx *gorm.DB) (err error) {
	var listUser = ListUser{ListID: l.ID, UserID: l.UserID, Creator: true}
	tx.Create(&listUser)
	return nil
}

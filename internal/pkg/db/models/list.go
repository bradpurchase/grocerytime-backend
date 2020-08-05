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
	DeletedAt *time.Time

	// Associations
	ListUsers    []ListUser
	Item         []Item
	GroceryTrips []GroceryTrip
}

// AfterCreate hook to automatically create some associated records
func (l *List) AfterCreate(tx *gorm.DB) (err error) {
	creator := true
	active := true
	listUser := ListUser{
		ListID:  l.ID,
		UserID:  l.UserID,
		Creator: &creator,
		Active:  &active,
	}
	if err := tx.Create(&listUser).Error; err != nil {
		return err
	}

	trip := GroceryTrip{ListID: l.ID, Name: "Trip 1"}
	if err := tx.Create(&trip).Error; err != nil {
		return err
	}

	return nil
}

package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email    string    `gorm:"type:varchar(100);unique_index;not null"`
	Password string    `gorm:"not null"`
	Name     string    `gorm:"type:varchar(100)"`

	LastSeenAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Associations
	Stores []Store
	Tokens []AuthToken
}

// AfterCreate hook to automatically create some associated records
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	// Create default store
	store := Store{UserID: u.ID, Name: "My Grocery Store"}
	if err := tx.Create(&store).Error; err != nil {
		return err
	}

	return
}

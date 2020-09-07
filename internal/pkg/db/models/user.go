package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Email     string    `gorm:"type:varchar(100);unique_index;not null"`
	Password  string    `gorm:"not null"`
	FirstName string    `gorm:"type:varchar(100)"`
	LastName  string    `gorm:"type:varchar(100)"`

	LastSeenAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Associations
	Stores []Store
	Tokens []AuthToken
}

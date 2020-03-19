package models

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Email      string    `gorm:"column:email;type:varchar(100);unique_index;not null"`
	Password   string    `gorm:"column:password;not null"`
	FirstName  string    `gorm:"column:first_name;type:varchar(100)"`
	LastName   string    `gorm:"column:last_name;type:varchar(100)"`
	LastSeenAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Associations
	Lists  []List
	Tokens []AuthToken
}
package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID                       uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email                    string     `gorm:"type:varchar(100);unique_index;not null"`
	Password                 string     `gorm:"not null"`
	Name                     string     `gorm:"type:varchar(100)"`
	PasswordResetToken       *uuid.UUID `gorm:"type:uuid"`
	PasswordResetTokenExpiry *time.Time

	LastSeenAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Associations
	Stores []Store
	Tokens []AuthToken
}

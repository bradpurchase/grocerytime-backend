package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID                       uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email                    string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password                 string     `gorm:"not null"`
	Name                     string     `gorm:"type:varchar(100)"`
	PasswordResetToken       *uuid.UUID `gorm:"type:uuid"`
	PasswordResetTokenExpiry *time.Time
	SiwaID                   *string `gorm:"type:varchar(255);uniqueIndex"`

	LastSeenAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Associations
	Stores []Store
	Tokens []AuthToken
}

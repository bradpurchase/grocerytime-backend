package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Device struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	Token  string    `gorm:"type:varchar(255);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	User User
}

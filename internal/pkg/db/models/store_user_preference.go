package models

import (
	"time"

	"gorm.io/gorm"

	uuid "github.com/satori/go.uuid"
)

type StoreUserPreference struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StoreUserID   uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	DefaultStore  bool      `gorm:"default:false;not null"`
	Notifications bool      `gorm:"default:true;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

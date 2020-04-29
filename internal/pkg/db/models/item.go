package models

import (
	"time"

	// Postgres dialect for GORM
	_ "github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

// Item defines the model for items
type Item struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ListID    uuid.UUID `gorm:"type:uuid;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Quantity  int       `gorm:"default:1;not null"`
	Completed bool      `gorm:"default:false;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations
	List List
}

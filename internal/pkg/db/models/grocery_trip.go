package models

import (
	"time"

	// Postgres dialect for GORM
	_ "github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

type GroceryTrip struct {
	ID   uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Name string    `gorm:"type:varchar(100)"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations
	List List
}

package models

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

type Item struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ListID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations
	List List
}

package models

import (
	uuid "github.com/satori/go.uuid"
)

// Category defines the model for item categories
type Category struct {
	ID   uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Name string    `gorm:"type:varchar(100);not null"`
}

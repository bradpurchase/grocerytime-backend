package models

import (
	"time"

	// Postgres dialect for GORM
	_ "github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

type ListUser struct {
	ID      uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ListID  uuid.UUID `gorm:"type:uuid;not null"`
	UserID  uuid.UUID `gorm:"type:uuid"`
	Email   string    `gorm:"type:varchar(100)"`
	Creator bool      `gorm:"default:false;not null"`
	Active  bool      `gorm:"default:true;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations
	List      List
	User      User
	ListUsers []ListUser
}

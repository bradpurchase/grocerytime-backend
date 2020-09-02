package models

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	// Random string generation for key/secret
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/utils"
)

type ApiClient struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(100);unique_index;not null"`
	Key       string    `gorm:"type:varchar(100);not null"`
	Secret    string    `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations
	Tokens []AuthToken
}

// BeforeCreate hook to generate key/secret
func (c *ApiClient) BeforeCreate(scope *gorm.Scope) (err error) {
	rand.Seed(time.Now().UnixNano())
	scope.SetColumn("Key", utils.RandString(24))
	scope.SetColumn("Secret", utils.RandString(24))
	return nil
}

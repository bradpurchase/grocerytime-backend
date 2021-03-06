package models

import (
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	// Random string generation for key/secret
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/utils"
)

type ApiClient struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Key       string    `gorm:"type:varchar(100);not null"`
	Secret    string    `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Associations
	Tokens []AuthToken `gorm:"foreignKey:ClientID"`
}

// BeforeCreate hook to generate key/secret
func (c *ApiClient) BeforeCreate(tx *gorm.DB) (err error) {
	rand.Seed(time.Now().UnixNano())
	c.Key = utils.RandString(24)
	c.Secret = utils.RandString(24)
	return
}

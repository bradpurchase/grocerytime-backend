package models

import (
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	// Random string generation for key/secret
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/utils"
)

type AuthToken struct {
	ID           uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ClientID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();not null"`
	UserID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();not null"`
	AccessToken  string    `gorm:"type:varchar(100);not null"`
	RefreshToken string    `gorm:"type:varchar(100);not null"`
	ExpiresIn    time.Time `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Associations
	Client ApiClient
	User   User
}

// BeforeCreate generates the AccessToken and RefreshToken, and sets
// ExpiresIn to 10 minutes from creation time so that access tokens frequently expire
func (c *AuthToken) BeforeCreate(tx *gorm.DB) (err error) {
	rand.Seed(time.Now().UnixNano())
	c.AccessToken = utils.RandString(20)
	c.RefreshToken = utils.RandString(20)
	c.ExpiresIn = time.Now().Add(time.Minute * 10)
	return
}

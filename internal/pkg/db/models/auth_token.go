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
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ClientID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();not null"`
	UserID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();not null"`
	AccessToken  string    `gorm:"type:varchar(100);not null;index:idx_auth_tokens_access_token"`
	RefreshToken string    `gorm:"type:varchar(100);not null"`
	ExpiresIn    time.Time `gorm:"not null"`
	DeviceName   string    `gorm:"type:varchar(100)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Associations
	Client ApiClient
	User   User
}

// BeforeCreate handles generating tokens and also handles old token cleanup
func (c *AuthToken) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate the AccessToken and RefreshToken, and set
	// ExpiresIn to 10 minutes from creation time so that access tokens
	// frequently expire.
	//
	// Note: refresh tokens are not currently in use. deprecate?
	rand.Seed(time.Now().UnixNano())
	c.AccessToken = utils.RandString(20)
	c.RefreshToken = utils.RandString(20)
	c.ExpiresIn = time.Now().Add(time.Minute * 10)

	// Clean up tokens after creation.
	//
	// Since our app is universal, it means a session can be held on both
	// iOS and iPadOS simultaneously. Therefore, we only delete a token if it's
	// replacing an existing token on the same device.
	authToken := &AuthToken{}
	tokenQuery := tx.
		Where("user_id = ? AND client_id = ? AND device_name = ?", c.UserID, c.ClientID, c.DeviceName).
		Delete(&authToken).
		Error
	if err := tokenQuery; err != nil {
		return err
	}
	return
}

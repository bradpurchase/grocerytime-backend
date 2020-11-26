package models

import (
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	// Random string generation for key/secret
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/utils"
)

// AuthTokensOffset is the number of tokens we retain when cleaning up
const AuthTokensOffset = 2

// AuthToken is a model that represents the auth_tokens table
type AuthToken struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
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

// BeforeCreate handles generating tokens and also handles old token cleanup
func (c *AuthToken) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate the AccessToken and RefreshToken, and sets
	// ExpiresIn to 10 minutes from creation time so that access tokens
	// frequently expire.
	//
	// Note: refresh tokens are not currently in use
	rand.Seed(time.Now().UnixNano())
	c.AccessToken = utils.RandString(20)
	c.RefreshToken = utils.RandString(20)
	c.ExpiresIn = time.Now().Add(time.Minute * 10)

	// Clean up old tokens after creation.
	// Since our app is universal, it means a session can be held on both
	// iOS and iPadOS simultaneously.
	//
	// Therefore, we only cleanup tokens if more than 2 exist.
	subquery := tx.
		Select("id").
		Table("auth_tokens").
		Where("user_id = ? AND client_id = ?", c.UserID, c.ClientID).
		Order("created_at DESC").
		Limit(100).
		Offset(AuthTokensOffset)
	tokenQuery := tx.
		Where("id IN (?)", subquery).
		Delete(&AuthToken{UserID: c.UserID, ClientID: c.ClientID}).
		Error
	if err := tokenQuery; err != nil {
		return err
	}

	return
}

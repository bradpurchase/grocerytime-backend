package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

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
func (c *AuthToken) BeforeCreate(scope *gorm.Scope) (err error) {
	scope.SetColumn("AccessToken", utils.RandString(20))
	scope.SetColumn("RefreshToken", utils.RandString(20))
	scope.SetColumn("ExpiresIn", time.Now().Add(time.Minute*10))

	return nil
}

package auth

import (
	"errors"
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a user account with an email and password
func CreateUser(db *gorm.DB, email string, password string, clientID uuid.UUID) (*models.User, error) {
	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Handle dupe email
	dupeUser := &models.User{}
	if err := db.Where("email = ?", email).First(&dupeUser).Error; !gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("An account with this email address already exists")
	}

	user := &models.User{
		Email:      email,
		Password:   string(passhash),
		LastSeenAt: time.Now(),
		Tokens: []models.AuthToken{
			{ClientID: clientID},
		},
		Lists: []models.List{
			{Name: "My Grocery List"},
		},
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Send an email upon user creation
	_, mailErr := mailer.SendNewUserEmail(user.Email)
	if mailErr != nil {
		return nil, mailErr
	}

	return user, nil
}

package auth

import (
	"errors"
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUser creates a user account with an email and password
func CreateUser(email string, password string, clientID uuid.UUID) (*models.User, error) {
	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Handle dupe email
	dupeUser := &models.User{}
	if err := db.Manager.Where("email = ?", email).First(&dupeUser).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("An account with this email address already exists")
	}

	user := &models.User{
		Email:      email,
		Password:   string(passhash),
		LastSeenAt: time.Now(),
		Tokens:     []models.AuthToken{{ClientID: clientID}},
	}
	if err := db.Manager.Create(&user).Error; err != nil {
		return nil, err
	}

	// Send an email upon user creation
	_, mailErr := mailer.SendNewUserEmail(user.Email)
	if mailErr != nil {
		return nil, mailErr
	}

	return user, nil
}

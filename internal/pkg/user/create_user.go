package user

import (
	"errors"
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a user account with details provided
func CreateUser(email string, password string, name string, clientID uuid.UUID) (*models.User, error) {
	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Check for dupe email
	var dupeCount int64
	if err := db.Manager.Model(&models.User{}).Where("email = ?", email).Count(&dupeCount).Error; err != nil {
		return nil, err
	}
	if dupeCount > 0 {
		return nil, errors.New("An account with this email address already exists")
	}

	user := &models.User{
		Name:       name,
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

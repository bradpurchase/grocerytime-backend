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
func CreateUser(email string, password string, name string, deviceName string, clientID uuid.UUID) (user models.User, err error) {
	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	// Check for dupe email
	var dupeCount int64
	if err := db.Manager.Model(&models.User{}).Where("email = ?", email).Count(&dupeCount).Error; err != nil {
		return user, err
	}
	if dupeCount > 0 {
		return user, errors.New("An account with this email address already exists")
	}

	user = models.User{
		Name:       name,
		Email:      email,
		Password:   string(passhash),
		LastSeenAt: time.Now(),
		Tokens:     []models.AuthToken{{ClientID: clientID, DeviceName: deviceName}},
	}
	if err := db.Manager.Create(&user).Error; err != nil {
		return user, err
	}

	// Send an email upon user creation
	_, mailErr := mailer.SendNewUserEmail(user.Email)
	if mailErr != nil {
		return user, mailErr
	}

	return user, nil
}

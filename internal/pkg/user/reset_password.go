package user

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	"golang.org/x/crypto/bcrypt"
)

// ResetPassword changes the password for a user after verifying the token
//
// The token here is a reset password token sent to a user when they request a
// password change as part of the forgot password flow
func ResetPassword(password string, token string) (resetUser *models.User, err error) {
	user := &models.User{}
	userQuery := db.Manager.
		Select("id, email").
		Where("password_reset_token = ? AND password_reset_token_expiry > now()", token).
		First(&user).
		Error
	if err := userQuery; err != nil {
		return resetUser, errors.New("token invalid or expired")
	}

	// Update the password and send an email to the user
	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return resetUser, err
	}
	if err := db.Manager.Model(&user).Update("password", string(passhash)).Error; err != nil {
		return resetUser, err
	}

	// Send the email to the user
	email := user.Email
	_, mailErr := mailer.SendPasswordResetEmail(email)
	if mailErr != nil {
		return resetUser, mailErr
	}

	return user, nil
}

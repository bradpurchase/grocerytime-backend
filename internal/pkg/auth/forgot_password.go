package auth

import (
	"errors"
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	uuid "github.com/satori/go.uuid"
)

// SendForgotPasswordEmail sends an email with a unique reset password link
// to an email address associated to a user
func SendForgotPasswordEmail(email string) (resetUser *models.User, err error) {
	user := &models.User{}
	if err := db.Manager.Where("email = ?", email).First(&user).Error; err != nil {
		return resetUser, err
	}

	// Generate a password reset token
	resetToken := uuid.NewV4()
	resetTokenExpiry := time.Now().Add(time.Minute * 10) // expires in 10 min
	updateQuery := db.Manager.
		Model(&user).
		Updates(models.User{PasswordResetToken: resetToken, PasswordResetTokenExpiry: resetTokenExpiry}).
		Error
	if err := updateQuery; err != nil {
		return resetUser, errors.New("could not generate reset token and expiry")
	}

	// Send the email to the user
	_, mailErr := mailer.SendForgotPasswordEmail(email, resetToken)
	if mailErr != nil {
		return resetUser, mailErr
	}

	return user, err
}

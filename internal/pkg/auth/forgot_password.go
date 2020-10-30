package auth

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// SendForgotPasswordEmail sends an email to an email address if it is associated
// with a user
func SendForgotPasswordEmail(email string) (user models.User, err error) {
	if err := db.Manager.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	return user, err
}

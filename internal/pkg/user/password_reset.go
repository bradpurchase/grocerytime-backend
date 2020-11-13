package user

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// VerifyPasswordResetToken retrieves a user associated with a password reset token
func VerifyPasswordResetToken(token interface{}) (tokenUser *models.User, err error) {
	user := &models.User{}
	query := db.Manager.
		Select("password_reset_token_expiry").
		Where("password_reset_token = ? AND password_reset_token_expiry > now()", token).
		First(&user).
		Error
	if err := query; err != nil {
		return tokenUser, errors.New("token expired")
	}
	return user, nil
}

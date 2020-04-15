package auth

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
)

// FetchAuthenticatedUser retrieves the user to satisfy AuthenticatedUserResolver
func FetchAuthenticatedUser(db *gorm.DB, header string) (interface{}, error) {
	token, err := RetrieveAccessToken(header)
	if err != nil {
		return nil, err
	}
	authToken := &models.AuthToken{}
	if err := db.Preload("User").Where("access_token = ? AND expires_in > now()", token).Last(&authToken).Error; err != nil {
		return nil, errors.New("Token is invalid or session has expired")
	}
	return authToken.User, nil
}

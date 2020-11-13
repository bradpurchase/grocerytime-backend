package auth

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// FetchAuthenticatedUser retrieves the user to satisfy AuthenticatedUserResolver
func FetchAuthenticatedUser(header string) (interface{}, error) {
	token, err := RetrieveAccessToken(header)
	if err != nil {
		return nil, err
	}
	authToken := &models.AuthToken{}
	query := db.Manager.
		Preload("User").
		Where("access_token = ?", token).
		Last(&authToken).
		Error
	if err := query; err != nil {
		return nil, errors.New("token invalid/expired")
	}
	return authToken.User, nil
}

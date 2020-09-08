package resolvers

import (
	"errors"

	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// TokenResolver fetches a token for a user authentication session. It requires
// a grantType argument which determines whether or not this is a login with a
// username and password in exchange for an access token, or if it's a refresh
// token being passed in exchange for a new access token.
//
// DEPRECATED: replaced by login mutation (LoginResolver)
func TokenResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	creds, err := auth.RetrieveClientCredentials(header.(string))
	if err != nil {
		return nil, err
	}
	apiClient := &models.ApiClient{}
	if err := db.Where("key = ? AND secret = ?", creds[0], creds[1]).First(&apiClient).Error; err != nil {
		return nil, err
	}

	switch grantType := p.Args["grantType"].(string); grantType {
	case "login":
		// In this case we accept an email and password, check that the email is in
		// the system and verify the password hash. Finally, we return a new
		// access token and refresh token after clearing any existing tokens
		email := p.Args["email"]
		password := p.Args["password"]
		if email == nil || password == nil {
			return nil, errors.New("Missing required arguments for login grant type: email, password")
		}

		//TODO i18n
		wrongCredsMsg := "Could not log you in with those details. Please try again!"
		user := &models.User{}
		if err := db.Where("email = ?", email.(string)).First(&user).Error; err != nil {
			return nil, errors.New(wrongCredsMsg)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password.(string))); err != nil {
			return nil, errors.New(wrongCredsMsg)
		}

		authToken := &models.AuthToken{UserID: user.ID, ClientID: apiClient.ID}
		if err := db.Where("user_id = ? AND client_id = ?", user.ID, apiClient.ID).Delete(&authToken).Error; err != nil {
			return nil, errors.New(wrongCredsMsg)
		}
		if err := db.Create(&authToken).Error; err != nil {
			return nil, errors.New(wrongCredsMsg)
		}

		return authToken, nil
	case "refreshToken":
		// In this case we are being asked for a new access token. In order for us
		// to grant the new token, we need to accept a refresh token and validate it.
		// Refresh tokens are valid if they belong to the user/client and are less
		// than a year old, as identified by the creation time of the AuthToken record
		refreshToken := p.Args["refreshToken"]
		if refreshToken == nil {
			return nil, errors.New("Missing required arguments for refreshToken grant type: refreshToken")
		}

		authToken := &models.AuthToken{}
		if err := db.Where("refresh_token = ? AND created_at >= now() - interval '1 year'", refreshToken.(string)).Last(&authToken).Error; err != nil {
			return nil, err
		}
		userID := authToken.UserID
		if err := db.Delete(&authToken).Error; err != nil {
			return nil, err
		}

		newAuthToken := &models.AuthToken{UserID: userID, ClientID: apiClient.ID}
		if err := db.Create(&newAuthToken).Error; err != nil {
			return nil, err
		}

		return newAuthToken, nil
	default:
		return nil, errors.New("Invalid value retrieved for grantType")
	}
}

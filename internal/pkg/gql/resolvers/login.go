package resolvers

import (
	"errors"

	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// LoginResolver fetches a token for an user authentication session
func LoginResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	creds, err := auth.RetrieveClientCredentials(header.(string))
	if err != nil {
		return nil, err
	}
	apiClient := &models.ApiClient{}
	if err := db.Manager.Where("key = ? AND secret = ?", creds[0], creds[1]).First(&apiClient).Error; err != nil {
		return nil, err
	}

	// In this case we accept an email and password, check that the email is in
	// the system and verify the password hash. Finally, we return a new
	// access token and refresh token after clearing any existing tokens
	email := p.Args["email"]
	password := p.Args["password"]
	if email == nil || password == nil {
		return nil, errors.New("missing required arguments for login: email, password")
	}

	//TODO i18n
	wrongCredsMsg := "Could not log you in with those details. Please try again!"
	var user models.User
	if err := db.Manager.Where("email = ?", email.(string)).First(&user).Error; err != nil {
		return nil, errors.New(wrongCredsMsg)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password.(string))); err != nil {
		return nil, errors.New(wrongCredsMsg)
	}

	var deviceName string
	if p.Args["deviceName"] != nil {
		deviceName = p.Args["deviceName"].(string)
	}
	authToken := &models.AuthToken{
		UserID:     user.ID,
		ClientID:   apiClient.ID,
		DeviceName: deviceName,
	}
	if err := db.Manager.Create(&authToken).Error; err != nil {
		return nil, errors.New(wrongCredsMsg)
	}

	return user, nil
}

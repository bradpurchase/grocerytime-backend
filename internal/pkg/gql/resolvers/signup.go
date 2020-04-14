package resolvers

import (
	"errors"

	"github.com/graphql-go/graphql"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// SignupResolver creates a new user account and assigns it a
func SignupResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	// Retrieve API client for the key/secret provided
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	creds, err := auth.RetrieveClientCredentials(header.(string))
	if err != nil {
		return nil, err
	}
	apiClient := &models.ApiClient{}
	if err := db.Where("key = ? AND secret = ?", creds[0], creds[1]).First(&apiClient).Error; err != nil {
		return nil, err
	}

	// Create a new user account with the args provided
	email := p.Args["email"].(string)
	password := p.Args["password"].(string)
	user, err := auth.CreateUser(db, email, password, apiClient.ID)
	if err != nil {
		return nil, errors.New("Hmm, we could not sign you up successfully. Please try again")
	}
	token := user.Tokens[0]
	return token, nil
}

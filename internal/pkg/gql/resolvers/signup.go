package resolvers

import (
	"github.com/graphql-go/graphql"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// SignupResolver creates a new user account and assigns it a
func SignupResolver(p graphql.ResolveParams) (interface{}, error) {
	// Retrieve API client for the key/secret provided
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	creds, err := auth.RetrieveClientCredentials(header.(string))
	if err != nil {
		return nil, err
	}
	apiClient := &models.ApiClient{}
	if err := db.Manager.Where("key = ? AND secret = ?", creds[0], creds[1]).First(&apiClient).Error; err != nil {
		return nil, err
	}

	// Create a new user account with the args provided
	email := p.Args["email"].(string)
	password := p.Args["password"].(string)
	name := p.Args["name"].(string)
	user, err := auth.CreateUser(email, password, name, apiClient.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

package resolvers

import (
	"github.com/graphql-go/graphql"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
)

// SignupResolver creates a new user account and assigns it a
func SignupResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()

	// Retrieve API client for the key/secret provided
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	_, err := auth.RetrieveClientCredentials(header.(string))
	if err != nil {
		return nil, err
	}

	// Create a new user account with the args provided
	email := p.Args["email"].(string)
	password := p.Args["password"].(string)
	user, err := auth.CreateUser(db, email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

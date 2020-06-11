package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

func SignupAndJoinListResolver(p graphql.ResolveParams) (interface{}, error) {
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
		return nil, err
	}

	// Check if a list with the listId provided exists and
	// create a list_users record for this user if so
	listID := p.Args["listId"].(uuid.UUID)
	if err := db.Where("id = ?", listID)
}

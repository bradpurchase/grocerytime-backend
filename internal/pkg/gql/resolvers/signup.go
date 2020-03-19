package resolvers

import (
	"time"

	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"

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
	password, err := bcrypt.GenerateFromPassword([]byte(p.Args["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Email:      p.Args["email"].(string),
		Password:   string(password),
		LastSeenAt: time.Now(),
		Tokens: []models.AuthToken{
			{ClientID: apiClient.ID},
		},
		Lists: []models.List{
			{Name: "My Grocery List"},
		},
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user.Tokens[0], nil
}

package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// SignInWithAppleResolver resolves the signInWithApple mutation
func SignInWithAppleResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	creds, err := auth.RetrieveClientCredentials(header.(string))
	if err != nil {
		return nil, err
	}
	var apiClient models.ApiClient
	if err := db.Manager.Where("key = ? AND secret = ?", creds[0], creds[1]).First(&apiClient).Error; err != nil {
		return nil, err
	}

	appScheme := p.Info.RootValue.(map[string]interface{})["App-Scheme"]

	identityToken := p.Args["identityToken"].(string)
	nonce := p.Args["nonce"].(string)
	email := p.Args["email"].(string)
	name := p.Args["name"].(string)
	user, err := auth.SignInWithApple(identityToken, nonce, email, name, appScheme.(string), apiClient.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

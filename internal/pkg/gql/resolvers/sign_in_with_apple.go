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
	apiClient := &models.ApiClient{}
	if err := db.Manager.Where("key = ? AND secret = ?", creds[0], creds[1]).First(&apiClient).Error; err != nil {
		return nil, err
	}

	// token, err := jwt.Parse(identityToken, func(t *jwt.Token) ([]byte, error) {
	// 	fmt.Println(t)
	// 	return nil, nil
	// })

	identityToken := p.Args["identityToken"].(string)
	email := p.Args["email"].(string)
	name := p.Args["name"].(string)
	user, err := auth.SignInWithApple(identityToken, email, name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

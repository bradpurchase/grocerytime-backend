package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/meals"
	"github.com/graphql-go/graphql"
)

// RecipesResolver returns Store records for the current user
func RecipesResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	userID := user.(models.User).ID
	recipes, err := meals.RetrieveRecipes(userID)
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

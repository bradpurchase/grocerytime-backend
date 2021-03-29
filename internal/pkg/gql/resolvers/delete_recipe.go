package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/meals"
	"github.com/graphql-go/graphql"
)

// DeleteRecipeResolver resolves the deleteRecipe mutation
func DeleteRecipeResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	recipeID := p.Args["id"]
	userID := user.ID
	recipe, err := meals.DeleteRecipe(recipeID, userID)
	if err != nil {
		return nil, err
	}

	return recipe, err
}

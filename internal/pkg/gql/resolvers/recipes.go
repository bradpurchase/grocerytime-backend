package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/meals"
	"github.com/graphql-go/graphql"
)

// RecipesResolver resolves the recipes query
func RecipesResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	userID := user.ID
	recipes, err := meals.RetrieveRecipes(userID, p.Args["mealType"])
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

// RecipeResolver resolves the recipe query
func RecipeResolver(p graphql.ResolveParams) (interface{}, error) {
	recipeID := p.Args["id"]
	recipe, err := meals.RetrieveRecipe(recipeID)
	if err != nil {
		return nil, err
	}
	return recipe, nil
}

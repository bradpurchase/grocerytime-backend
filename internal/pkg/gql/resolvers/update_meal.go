package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/meals"
	"github.com/graphql-go/graphql"
)

// UpdateMealResolver resolves the updateMeal mutation
func UpdateMealResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	_, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	appScheme := p.Info.RootValue.(map[string]interface{})["App-Scheme"]
	meal, err := meals.UpdateMeal(p.Args, appScheme.(string))
	if err != nil {
		return nil, err
	}

	return meal, err
}

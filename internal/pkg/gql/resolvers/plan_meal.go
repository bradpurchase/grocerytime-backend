package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/meals"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/notifications"
	"github.com/graphql-go/graphql"
)

// PlanMealResolver resolves the planMeal mutation
func PlanMealResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	meal, err := meals.PlanMeal(user.ID, p.Args)
	if err != nil {
		return nil, err
	}

	appScheme := p.Info.RootValue.(map[string]interface{})["App-Scheme"]
	if appScheme == nil {
		return meal, err
	}
	go notifications.MealPlanned(meal, appScheme.(string))

	return meal, nil
}

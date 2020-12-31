package resolvers

import (
	"time"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/meals"
	"github.com/graphql-go/graphql"
)

// PlannedMealsResolver resolves the plannedMeals query
func PlannedMealsResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	year := p.Args["year"]
	weekNumber := p.Args["weekNumber"]
	// we allow no value for year and weekNumber args and use current time for this case
	if year == nil || weekNumber == nil {
		t := time.Now()
		year, weekNumber = t.ISOWeek()
	}
	meals, err := meals.PlannedMeals(user.ID, weekNumber.(int), year.(int))
	if err != nil {
		return nil, err
	}
	return meals, nil
}

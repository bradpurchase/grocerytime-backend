package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/meals"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// DeleteMealResolver resolves the deleteMeal mutation
func DeleteMealResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	mealID := p.Args["id"].(uuid.UUID)
	userID := user.ID
	appScheme := p.Info.RootValue.(map[string]interface{})["App-Scheme"]
	meal, err := meals.DeleteMeal(mealID, userID, appScheme.(string))
	if err != nil {
		return nil, err
	}

	return meal, err
}

package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/meals"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// MealResolver resolves the meals query
func MealResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	mealID := p.Args["id"].(uuid.UUID)
	meal, err := meals.RetrieveMealForUser(mealID, user.ID)
	if err != nil {
		return nil, err
	}
	return meal, nil
}

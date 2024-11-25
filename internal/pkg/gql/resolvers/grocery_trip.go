package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// GroceryTripResolver retrieves a grocery trip by ID
func GroceryTripResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	_, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	tripID := p.Args["id"].(uuid.UUID)
	trip, err := trips.RetrieveTrip(tripID)
	if err != nil {
		return nil, err
	}
	return trip, nil
}

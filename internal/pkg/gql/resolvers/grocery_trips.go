package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

// GroceryTripsResolver resolves the trips mutation
func GroceryTripsResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	storeID := p.Args["storeId"]
	userID := user.(models.User).ID
	completed := p.Args["completed"].(bool)
	trips, err := trips.RetrieveTrips(storeID, userID, completed)
	if err != nil {
		return nil, err
	}
	return trips, nil
}

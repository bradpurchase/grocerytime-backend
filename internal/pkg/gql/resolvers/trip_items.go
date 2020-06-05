package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

// TripItemsResolver resolves the items query by retrieving the items in a trip
func TripItemsResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	tripID := p.Source.(models.GroceryTrip).ID
	items, err := trips.RetrieveItems(db, tripID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// TripItemsCountResolver returns a count of items in a trip
func TripItemsCountResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	tripID := p.Source.(models.GroceryTrip).ID
	items, err := trips.RetrieveItems(db, tripID)
	if err != nil {
		return 0, err
	}
	itemsCount := len(items.([]models.Item))
	return itemsCount, nil
}

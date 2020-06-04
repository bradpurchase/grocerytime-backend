package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

// GroceryTripResolver retrieves the current grocery trip in a list
// as part of the List query
func GroceryTripResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	listID := p.Source.(*models.List).ID
	userID := user.(models.User).ID
	trip, err := trips.RetrieveCurrentTripInList(db, listID, userID)
	if err != nil {
		return nil, err
	}
	return trip, nil
}

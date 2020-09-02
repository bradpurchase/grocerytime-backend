package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

// ReorderItemResolver updates the position of an item with the provided params
func ReorderItemResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	_, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	itemID := p.Args["itemId"]
	position := p.Args["position"].(int)
	trip, err := trips.ReorderItem(db, itemID, position)
	if err != nil {
		return nil, err
	}
	return trip, err
}

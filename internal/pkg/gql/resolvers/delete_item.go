package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

// DeleteItemResolver deletes an item by itemId param
func DeleteItemResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	_, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	itemID := p.Args["itemId"]
	item, err := trips.DeleteItem(db, itemID)
	if err != nil {
		return nil, err
	}
	return item, err
}

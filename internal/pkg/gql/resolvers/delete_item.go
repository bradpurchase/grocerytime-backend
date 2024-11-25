package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// DeleteItemResolver deletes an item by itemId param
func DeleteItemResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	_, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	itemID := p.Args["itemId"].(uuid.UUID)
	item, err := trips.DeleteItem(itemID)
	if err != nil {
		return nil, err
	}
	return item, err
}

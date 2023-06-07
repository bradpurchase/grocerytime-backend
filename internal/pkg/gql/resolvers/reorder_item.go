package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// ReorderItemResolver updates the position of an item with the provided params
func ReorderItemResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	_, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	itemID := p.Args["itemId"].(uuid.UUID)
	position := p.Args["position"].(int)
	trip, err := trips.ReorderItem(itemID, position)
	if err != nil {
		return nil, err
	}
	return trip, err
}

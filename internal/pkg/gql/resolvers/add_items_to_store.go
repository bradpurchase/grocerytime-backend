package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

// AddItemsToStore resolves the addItemsToStore mutation
func AddItemsToStore(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	userID := user.ID
	items, err := trips.AddItemsToStore(userID, p.Args)
	if err != nil {
		return nil, err
	}
	return items, nil
}

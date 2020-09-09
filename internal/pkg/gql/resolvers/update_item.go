package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

// UpdateItemResolver updates the properties of an item with the provided params
func UpdateItemResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	_, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	item, err := trips.UpdateItem(p.Args)
	if err != nil {
		return nil, err
	}
	return item, err
}

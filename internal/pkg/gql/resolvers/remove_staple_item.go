package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// RemoveStapleItem resolves the removeStapleItem mutation
func RemoveStapleItem(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	_, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	itemID, err := uuid.FromString(p.Args["itemId"].(string))
	if err != nil {
		return nil, err
	}
	item, err := stores.RemoveStapleItem(itemID)
	if err != nil {
		return nil, err
	}
	return item, err
}

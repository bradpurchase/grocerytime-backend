package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// DeleteStoreResolver resolves the deleteStore mutation by deleting a store
// and its associated store users and items
func DeleteStoreResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	storeID := p.Args["storeId"].(uuid.UUID)
	store, err := stores.DeleteStore(storeID, user.ID)
	if err != nil {
		return nil, err
	}
	return store, nil
}

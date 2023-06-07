package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// StoreResolver resolves the store GraphQL query by retrieving a store by ID param
func StoreResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	storeID := p.Args["id"].(uuid.UUID)
	store, err := stores.RetrieveStoreForUser(storeID, user.ID)
	if err != nil {
		return nil, err
	}
	return store, nil
}

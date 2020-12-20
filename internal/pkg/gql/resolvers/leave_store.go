package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
)

// LeaveStoreResolver resolves the leaveStore resolver by removing the current user from the store
func LeaveStoreResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	// Verify that the store with the ID provided exists
	storeID := p.Args["storeId"]
	storeUser, err := stores.RemoveUserFromStore(user, storeID)
	if err != nil {
		return nil, err
	}
	return storeUser, nil
}

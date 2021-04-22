package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"

	"github.com/graphql-go/graphql"
)

// JoinStoreResolver adds the current user to a store properly by removing
// the email and replacing it with the user ID
func JoinStoreResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	storeID := p.Args["storeId"]
	storeUser, err := stores.AddUserToStore(user, storeID)
	if err != nil {
		return nil, err
	}
	return storeUser, nil
}

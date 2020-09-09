package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
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

	// Verify that the store with the ID provided exists
	storeID := p.Args["storeId"]
	storeUser, err := stores.AddUserToStore(user.(models.User), storeID)
	if err != nil {
		return nil, err
	}
	return storeUser, nil
}

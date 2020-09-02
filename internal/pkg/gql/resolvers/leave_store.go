package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
)

// LeaveStoreResolver resolves the leaveStore resolver by removing the current user from the store
func LeaveStoreResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	// Verify that the store with the ID provided exists
	storeID := p.Args["storeId"]
	storeUser, err := stores.RemoveUserFromStore(db, user.(models.User), storeID)
	if err != nil {
		return nil, err
	}
	return storeUser, nil
}

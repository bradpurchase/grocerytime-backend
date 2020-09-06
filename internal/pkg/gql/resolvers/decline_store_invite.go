package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
)

// DeclineStoreInviteResolver resolves the declineStoreResolver resolver by calling
// stores.RemoveUserFromStore function which handles removing the StoreUser record
// and emailing the store creator about the invite being declined
func DeclineStoreInviteResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	storeID := p.Args["storeId"]
	storeUser, err := stores.RemoveUserFromStore(db, user.(models.User), storeID)
	if err != nil {
		return nil, err
	}
	return storeUser, nil
}
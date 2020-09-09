package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
)

// DeclineStoreInviteResolver resolves the declineStoreResolver resolver by calling
// stores.RemoveUserFromStore function which handles removing the StoreUser record
// and emailing the store creator about the invite being declined
func DeclineStoreInviteResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	storeID := p.Args["storeId"]
	storeUser, err := stores.RemoveUserFromStore(user.(models.User), storeID)
	if err != nil {
		return nil, err
	}
	return storeUser, nil
}

package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
)

// UpdateStoreUserPrefsResolver resolves the updateStoreUserPrefs mutation
func UpdateStoreUserPrefsResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	// Find the StoreUser record from the storeId arg provided and current user ID
	storeID := p.Args["storeId"]
	userID := user.(models.User).ID
	storeUserID, err := stores.RetrieveStoreUserID(storeID, userID)
	if err != nil {
		return nil, err
	}
	storeUserPrefs, err := stores.UpdateStoreUserPrefs(storeUserID, p.Args)
	if err != nil {
		return nil, err
	}
	return storeUserPrefs, nil
}

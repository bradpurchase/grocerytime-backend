package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
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
	var storeUser models.StoreUser
	storeUserQuery := db.Manager.
		Select("id").
		Where("store_id = ? AND user_id = ?", storeID, userID).
		Find(&storeUser).
		Error
	if err := storeUserQuery; err != nil {
		return nil, err
	}
	storeUserID := storeUser.ID
	storeUserPrefs, err := stores.UpdateStoreUserPrefs(storeUserID, p.Args)
	if err != nil {
		return nil, err
	}
	return storeUserPrefs, nil
}

package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"

	"github.com/graphql-go/graphql"
)

// JoinStoreWithShareCodeResolver resolves the joinStoreWithShareCode mutation
func JoinStoreWithShareCodeResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}
	appScheme := p.Info.RootValue.(map[string]interface{})["App-Scheme"]

	// Verify that the store with the ID provided exists
	storeID := p.Args["storeId"]
	var store models.Store
	if err := db.Manager.Model(&models.Store{}).Where("id = ?", storeID).First(&store).Error; err != nil {
		return nil, err
	}

	code := p.Args["code"].(string)
	storeUser, err := stores.AddUserToStoreWithCode(user, store, code, appScheme.(string))
	if err != nil {
		return nil, err
	}
	return storeUser, nil
}

package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
)

// UpdateStoreResolver resolves the updateStore mutation by updating the properties of a store
func UpdateStoreResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	store, err := stores.UpdateStoreForUser(db, user.(models.User).ID, p.Args)
	if err != nil {
		return nil, err
	}

	return store, nil
}

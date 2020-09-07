package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
)

// CreateStoreResolver creates a new store for the currently authenticated user
func CreateStoreResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	userID := user.(models.User).ID
	store, err := stores.CreateStore(db, userID, p.Args["name"].(string))
	if err != nil {
		return nil, err
	}
	return store, nil
}

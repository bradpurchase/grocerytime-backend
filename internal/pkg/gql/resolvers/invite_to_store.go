package resolvers

import (
	"errors"
	"strings"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
)

// InviteToStoreResolver resolves the inviteToStore mutation by creating a pending
// store_users record for the given storeId and email
func InviteToStoreResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}
	userEmail := user.Email

	storeID := p.Args["storeId"]
	invitedUserEmail := strings.TrimSpace(p.Args["email"].(string))
	if userEmail == invitedUserEmail {
		return models.StoreUser{}, errors.New("cannot invite yourself to a store")
	}
	storeUser, err := stores.InviteToStoreByEmail(storeID, invitedUserEmail)
	if err != nil {
		return nil, err
	}
	return storeUser, nil
}

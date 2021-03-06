package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	"github.com/graphql-go/graphql"
)

// InvitedStoresResolver resolves the invitedStores query by retrieving stores
// that the current user has been invited to
func InvitedStoresResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	invites, err := stores.RetrieveInvitedUserStores(user)
	if err != nil {
		return nil, err
	}

	return invites, nil
}

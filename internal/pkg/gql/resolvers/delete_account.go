package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/user"
	"github.com/graphql-go/graphql"
)

// DeleteAccountResolver resolves the deleteAccount mutation
func DeleteAccountResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	authUser, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	deletedUser, err := user.DeleteAccount(authUser)
	if err != nil {
		return nil, err
	}
	return deletedUser, nil
}

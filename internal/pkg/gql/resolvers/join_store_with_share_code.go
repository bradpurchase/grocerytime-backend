package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
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

	code := p.Args["code"].(string)
	storeUser, err := stores.AddUserToStoreWithCode(user, code, appScheme.(string))
	if err != nil {
		return nil, err
	}
	return storeUser, nil
}

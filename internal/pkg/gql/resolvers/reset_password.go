package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/user"
	"github.com/graphql-go/graphql"
)

// ResetPasswordResolver resolves the resetPassword mutation
func ResetPasswordResolver(p graphql.ResolveParams) (interface{}, error) {
	password := p.Args["password"].(string)
	token := p.Args["token"].(string)
	user, err := user.ResetPassword(password, token)
	if err != nil {
		return nil, err
	}
	return user, nil
}

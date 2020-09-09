package resolvers

import (
	"github.com/graphql-go/graphql"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// AuthenticatedUserResolver resolves me GraphQL query by returning the
// authenticated user, or an error if no authenticated user exists
func AuthenticatedUserResolver(p graphql.ResolveParams) (interface{}, error) {

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}
	return user, nil
}

// BasicUserResolver resolves the creator field in the StoreType by retrieving
// basic information about a user (email, first name, last name)
func BasicUserResolver(p graphql.ResolveParams) (interface{}, error) {
	user := &models.User{}
	if err := db.Manager.Where("id = ?", p.Source.(models.Store).UserID).Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

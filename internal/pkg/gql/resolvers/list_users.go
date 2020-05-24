package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"
	"github.com/graphql-go/graphql"
)

// ListUsersResolver resolves the listUsers field in the lists query
// by retrieving the list users belonging to a list
func ListUsersResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	listID := p.Source.(*models.List).ID
	listUsers, err := grocerylist.RetrieveListUsers(db, listID)
	if err != nil {
		return nil, err
	}
	return listUsers, nil
}

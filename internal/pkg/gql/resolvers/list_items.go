package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"
	"github.com/graphql-go/graphql"
)

// ListItemsResolver resolves the items query by retrieving the items in a list
func ListItemsResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	listID := p.Source.(*models.List).ID
	items, err := grocerylist.RetrieveItemsInList(db, listID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

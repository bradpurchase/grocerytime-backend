package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"
	"github.com/graphql-go/graphql"
)

// ListsResolver returns List records for the current user
func ListsResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	userID := p.Source.(models.User).ID
	lists, err := grocerylist.RetrieveUserLists(db, userID)
	if err != nil {
		return nil, err
	}

	return lists, nil
}

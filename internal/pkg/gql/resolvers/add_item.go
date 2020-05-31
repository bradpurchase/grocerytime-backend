package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"
	"github.com/graphql-go/graphql"
)

// AddItemResolver adds an item to the list with the listId param
func AddItemResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	item, err := grocerylist.AddItemToList(db, user.(models.User).ID, p.Args)
	if err != nil {
		return nil, err
	}
	return item, nil
}

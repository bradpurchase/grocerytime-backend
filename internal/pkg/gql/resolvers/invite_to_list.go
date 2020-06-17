package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"
	"github.com/graphql-go/graphql"
)

// InviteToList resolves the inviteToList mutation by creating a pending
// list_users record for the given listId and email
func InviteToListResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	_, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	// Verify that the list with the ID provided exists
	listID := p.Args["listId"]
	email := p.Args["email"].(string)
	listUser, err := grocerylist.InviteToListByEmail(db, listID, email)
	if err != nil {
		return nil, err
	}
	return listUser, nil
}

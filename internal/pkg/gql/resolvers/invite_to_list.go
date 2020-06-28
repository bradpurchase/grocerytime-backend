package resolvers

import (
	"errors"
	"strings"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"
	"github.com/graphql-go/graphql"
)

// InviteToListResolver resolves the inviteToList mutation by creating a pending
// list_users record for the given listId and email
func InviteToListResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}
	userEmail := user.(models.User).Email

	// Verify that the list with the ID provided exists
	listID := p.Args["listId"]
	invitedUserEmail := strings.TrimSpace(p.Args["email"].(string))
	if userEmail != invitedUserEmail {
		return models.ListUser{}, errors.New("cannot invite yourself to this list")
	}
	listUser, err := grocerylist.InviteToListByEmail(db, listID, invitedUserEmail)
	if err != nil {
		return nil, err
	}
	return listUser, nil
}

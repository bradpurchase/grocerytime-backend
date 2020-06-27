package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"
	"github.com/graphql-go/graphql"
)

// DeclineListInviteResolver resolves the declineListInvite resolver by calling
// grocerylist.RemoveUserFromList function which handles removing the ListUser record
// and emailing the list creator about the invite being declined
func DeclineListInviteResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	listID := p.Args["listId"]
	listUser, err := grocerylist.RemoveUserFromList(db, user.(models.User), listID)
	if err != nil {
		return nil, err
	}
	return listUser, nil
}

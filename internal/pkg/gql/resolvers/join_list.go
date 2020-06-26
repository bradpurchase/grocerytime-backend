package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"

	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

// JoinListResolver adds the current user to a list properly by removing
// the email and replacing it with the user ID
func JoinListResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	// Verify that the list with the ID provided exists
	listID := p.Args["listId"].(uuid.UUID)
	listUser, err := grocerylist.AddUserToList(db, user.(models.User), listID)
	if err != nil {
		return nil, err
	}
	return listUser, nil
}

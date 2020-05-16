package resolvers

import (
	"github.com/graphql-go/graphql"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"
)

// AddListUserResolver adds a ListUser to a List. It requires a listId variable,
// corresponding with a List, and an email variable.
//
// If the email address doesn't correspond with a User account, it creates one
// and sends an email the new user to complete signup.
func AddListUserResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}
	userID := user.(models.User).ID

	// Verify that the list with the ID provided exists and belongs to the current user
	list := &models.List{}
	if err := db.Where("id = ? AND user_id = ?", p.Args["listId"], userID).First(&list).Error; err != nil {
		return nil, err
	}

	sharedUserID := p.Args["userId"].(string)
	listUser, err := grocerylist.AddUserToList(db, sharedUserID, list)
	if err != nil {
		return nil, err
	}
	return listUser, nil
}

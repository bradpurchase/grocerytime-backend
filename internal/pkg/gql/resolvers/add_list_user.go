package resolvers

import (
	"errors"
	"log"

	"github.com/graphql-go/graphql"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
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

	// Verify that the list with the ID provided exists
	listID := p.Args["listId"]
	if err := db.Where("id = ?", listID).First(&models.List{}).Error; err != nil {
		return nil, errors.New("list not found")
	}

	// Verify that the current user owns this list and can add other users
	listCreatorUserID := user.(models.User).ID
	if err := db.Where("user_id = ? AND creator = ?", listCreatorUserID, true).First(&models.ListUser{}).Error; err != nil {
		return nil, errors.New("user is not the creator of this list and cannot add other users to it")
	}

	// Check if the email address provided corresponds with an existing user
	email := p.Args["email"].(string)
	if db.Where("email = ?", email).First(&models.User{}).RecordNotFound() {
		//TODO User doesn't exist; need to create 'em
		log.Println("user doesnt exist")
	} else {
		//TODO User exists.. we can just add them to the list
		log.Println("user exists")
	}

	return nil, nil
}

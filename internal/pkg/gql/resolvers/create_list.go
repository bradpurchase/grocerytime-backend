package resolvers

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"
	"github.com/graphql-go/graphql"
)

// CreateListResolver resolves a GraphQL mutation for creating a new List
// for the currently authenticated user
func CreateListResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	//TODO this needs to be moved into the grocerylist package
	// Check if this user has a list with this name already
	userID := user.(models.User).ID
	listName := p.Args["name"].(string)
	dupeList, _ := grocerylist.RetrieveListForUserByName(db, listName, userID)
	if dupeList.Name != "" {
		return nil, errors.New("You already have a list with this name")
	}

	list := &models.List{UserID: userID, Name: listName}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

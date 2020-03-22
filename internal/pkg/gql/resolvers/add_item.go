package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/graphql-go/graphql"
)

// AddItemResolver adds an item to the list with the listId variable
func AddItemResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	// Verify that the current user belongs in this list
	userID := user.(models.User).ID
	listUser := &models.ListUser{}
	if err := db.Where("list_id = ? AND user_id = ?", p.Args["listId"], userID).First(&listUser).Error; err != nil {
		return nil, err
	}

	item := &models.Item{
		ListID:   listUser.ListID,
		UserID:   userID,
		Name:     p.Args["name"].(string),
		Quantity: p.Args["quantity"].(int),
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}

	return item, nil
}

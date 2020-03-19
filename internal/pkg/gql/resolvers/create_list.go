package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
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

	//TODO handle case where user has a list with this name already
	list := &models.List{
		UserID: user.(models.User).ID,
		Name:   p.Args["name"].(string),
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

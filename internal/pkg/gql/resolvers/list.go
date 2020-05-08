package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/grocerylist"
	"github.com/graphql-go/graphql"
)

// ListResolver resolves the list GraphQL query by retrieving a list by ID param
func ListResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(db, header.(string))
	if err != nil {
		return nil, err
	}

	list, err := grocerylist.RetrieveListForUser(db, p.Args["id"], user.(models.User).ID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// SharableListResolver resolves the sharableList GraphQL query by retrieving
// basic info about a list for sharing purposes
func SharableListResolver(p graphql.ResolveParams) (interface{}, error) {
	db := db.FetchConnection()
	defer db.Close()

	list, err := grocerylist.RetrieveSharableList(db, p.Args["id"])
	if err != nil {
		return nil, err
	}
	return list, nil
}

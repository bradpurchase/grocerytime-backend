package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

// ItemSearchResolver resolves the itemSearch query
func ItemSearchResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	name := p.Args["name"].(string)
	userID := user.(models.User).ID
	item, err := trips.SearchForItemByName(name, userID)
	if err != nil {
		return nil, err
	}
	return item, err
}

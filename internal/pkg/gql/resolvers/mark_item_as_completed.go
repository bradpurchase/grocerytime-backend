package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

// MarkItemAsCompletedResolver resolves the markItemAsCompleted mutation
func MarkItemAsCompletedResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	userID := user.(models.User).ID
	name := p.Args["name"].(string)
	item, err := trips.MarkItemAsCompleted(name, userID)
	if err != nil {
		return nil, err
	}
	return item, err
}

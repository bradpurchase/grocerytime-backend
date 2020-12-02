package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/notifications"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	"github.com/graphql-go/graphql"
)

// AddItemToTrip resolves the addItemToTrip mutation
func AddItemToTrip(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	userID := user.(models.User).ID
	item, err := trips.AddItem(userID, p.Args)
	if err != nil {
		return nil, err
	}

	// Get the app scheme (i.e. Debug, Beta, Release) as we need to pass it to
	// the notifications package so it can use the proper apns certificate
	appScheme := p.Info.RootValue.(map[string]interface{})["App-Scheme"]
	if appScheme != nil {
		go notifications.ItemAdded(item, appScheme.(string))
	}

	return item, nil
}

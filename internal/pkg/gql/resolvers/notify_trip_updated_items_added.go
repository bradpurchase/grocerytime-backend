package resolvers

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/notifications"
	"github.com/graphql-go/graphql"
)

// NotifyTripUpdatedItemsAddedResolver resolves the notifyTripUpdatedItemsAdded mutation
func NotifyTripUpdatedItemsAddedResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return false, err
	}

	// Get the app scheme (i.e. Debug, Beta, Release) as we need to pass it to
	// the notifications package so it can use the proper apns certificate
	appScheme := p.Info.RootValue.(map[string]interface{})["App-Scheme"]
	if appScheme != nil {
		userID := user.ID
		storeID := p.Args["storeId"]
		numItemsAdded := p.Args["numItemsAdded"].(int)
		go notifications.ItemsAdded(userID, storeID, numItemsAdded, appScheme.(string))
		return true, nil
	}
	return false, errors.New("no app scheme provided")
}

package resolvers

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/notifications"
	"github.com/graphql-go/graphql"
)

// AddDeviceResolver creates a new user account and assigns it a
func AddDeviceResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	token := p.Args["token"].(string)
	userID := user.ID
	device, err := notifications.StoreDeviceToken(token, userID)
	if err != nil {
		return nil, err
	}
	return device, nil
}

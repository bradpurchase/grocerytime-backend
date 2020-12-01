package resolvers

import (
	"errors"
	"fmt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/auth"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/notifications"
	"github.com/graphql-go/graphql"
)

type PushNotification struct {
	ApnsID string
}

// SendPushNotificationResolver resolves the sendPushNotification mutation
func SendPushNotificationResolver(p graphql.ResolveParams) (interface{}, error) {
	header := p.Info.RootValue.(map[string]interface{})["Authorization"]
	user, err := auth.FetchAuthenticatedUser(header.(string))
	if err != nil {
		return nil, err
	}

	title := p.Args["title"].(string)
	body := p.Args["body"].(string)
	userID := user.(models.User).ID
	deviceTokens, err := notifications.DeviceTokensForUser(userID)
	if err != nil {
		return nil, errors.New("could not fetch device tokens for user")
	}
	var pushNotifications []PushNotification
	for i := range deviceTokens {
		apnsID, err := notifications.Send(title, body, deviceTokens[i])
		if err != nil {
			return nil, fmt.Errorf("could not send notification at index %d", i)
		}
		pushNotifications = append(pushNotifications, PushNotification{ApnsID: apnsID})
	}
	return pushNotifications, nil
}

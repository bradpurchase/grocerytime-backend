package notifications

import (
	"fmt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// ItemAdded sends a push notification to store users about a new item
func ItemAdded(item *models.Item, appScheme string) (err error) {
	title := "Trip Updated"
	body := fmt.Sprintf("%v added to your Backend trip", item.Name)
	token := "e9c5c8cc94425b19a6a0126608fcb9e1ea5455101db2d79959e56b9305bc1f41"
	_, e := Send(title, body, token, appScheme)
	if e != nil {
		return err
	}
	return nil
}

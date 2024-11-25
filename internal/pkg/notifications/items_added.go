package notifications

import (
	"fmt"
	"log"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// ItemsAdded sends a push notification to store users about a new item
func ItemsAdded(userID uuid.UUID, storeID uuid.UUID, numItemsAdded int, appScheme string) {
	var store models.Store
	if err := db.Manager.Select("id, name").Where("id = ?", storeID).First(&store).Error; err != nil {
		log.Println(err)
	}

	deviceTokens, err := StoreUserTokens(store.ID, userID)
	if err != nil {
		log.Println(err)
	}

	title := "Trip Updated"
	body := fmt.Sprintf("%d items added to your %v trip", numItemsAdded, store.Name)
	if numItemsAdded == 1 {
		body = fmt.Sprintf("An item was added to your %v trip", store.Name)
	}
	for i := range deviceTokens {
		storeIDStr := storeID.String()
		Send(title, body, deviceTokens[i], "Store", storeIDStr, appScheme)
	}
}

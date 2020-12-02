package notifications

import (
	"fmt"
	"log"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// ItemAdded sends a push notification to store users about a new item
func ItemAdded(item *models.Item, appScheme string) {
	// Fetch the store users excluding the one who created this item
	var store models.Store
	groceryTripQuery := db.Manager.
		Select("stores.id, stores.name").
		Joins("INNER JOIN grocery_trips ON grocery_trips.store_id = stores.id").
		Where("grocery_trips.id = ?", item.GroceryTripID).
		First(&store).
		Error
	if err := groceryTripQuery; err != nil {
		log.Println(err)
	}

	storeID := store.ID
	deviceTokens, err := StoreUserTokens(storeID, item)
	if err != nil {
		log.Println(err)
	}

	title := "Trip Updated"
	body := fmt.Sprintf("%v added to your %v trip", item.Name, store.Name)
	for i := range deviceTokens {
		Send(title, body, deviceTokens[i], appScheme)
	}
}

// StoreUserTokens fetches apns device tokens for all store users associated
// with the item provided, excluding for those users who have disabled
// notifications in store user preference settings
func StoreUserTokens(storeID uuid.UUID, item *models.Item) (tokens []string, err error) {
	var storeUsers []models.StoreUser
	storeUsersQuery := db.Manager.
		Select("store_users.user_id").
		Joins("INNER JOIN store_user_preferences ON store_user_preferences.store_user_id = store_users.id").
		Where("store_users.store_id = ?", storeID).
		Where("store_users.user_id NOT IN (?)", item.UserID).
		Where("store_user_preferences.notifications = ?", true).
		Find(&storeUsers).
		Error
	if err := storeUsersQuery; err != nil {
		return tokens, err
	}

	// Fetch the tokens for each store user
	var t []string
	for i := range storeUsers {
		userTokens, err := DeviceTokensForUser(storeUsers[i].UserID)
		if err != nil {
			return tokens, err
		}
		t = append(t, userTokens...)
	}
	return t, nil
}

package notifications

import (
	"fmt"
	"log"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// ItemsAdded sends a push notification to store users about a new item
func ItemsAdded(userID uuid.UUID, storeID interface{}, numItemsAdded int, appScheme string) {
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
		Send(title, body, deviceTokens[i], storeID.(string), appScheme)
	}
}

// StoreUserTokens fetches apns device tokens for all store users associated
// with the item provided, excluding for those users who have disabled
// notifications in store user preference settings
func StoreUserTokens(storeID uuid.UUID, userID uuid.UUID) (tokens []string, err error) {
	var storeUsers []models.StoreUser
	storeUsersQuery := db.Manager.
		Select("store_users.user_id").
		Joins("INNER JOIN store_user_preferences ON store_user_preferences.store_user_id = store_users.id").
		Where("store_users.store_id = ?", storeID).
		Where("store_users.user_id NOT IN (?)", userID).
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

package notifications

import (
	"fmt"
	"log"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// UserJoinedStore sends a push notification when a user joins a store via share code
func UserJoinedStore(user models.User, storeID interface{}, appScheme string) {
	var store models.Store
	if err := db.Manager.Select("id, name").Where("id = ?", storeID).First(&store).Error; err != nil {
		log.Println(err)
	}

	deviceTokens, err := StoreUserTokens(store.ID, user.ID)
	if err != nil {
		log.Println(err)
	}

	title := fmt.Sprintf("%v Joined Your Store", user.Name)
	body := fmt.Sprintf("%v has just joined your %v store. Now you can plan groceries and meals together!", user.Name, store.Name)
	for i := range deviceTokens {
		Send(title, body, deviceTokens[i], "Store", storeID.(string), appScheme)
	}
}

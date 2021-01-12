package notifications

import (
	"fmt"
	"log"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// MealPlanned sends a push notification to store users about a new item
func MealPlanned(meal *models.Meal, appScheme string) {
	var user models.User
	userQuery := db.Manager.
		Select("name").
		Where("id = ?", meal.UserID).
		Last(&user).
		Error
	if err := userQuery; err != nil {
		log.Println(err)
	}

	title := "New Meal Planned"
	body := fmt.Sprintf("%v added to your meal plan", user.Name)
	//TODO: should this be MealUserTokens? send to meal users instead of store users
	deviceTokens, err := StoreUserTokens(meal.StoreID, meal.UserID)
	if err != nil {
		log.Println(err)
	}
	for i := range deviceTokens {
		Send(title, body, deviceTokens[i], meal.ID.String(), appScheme)
	}
}

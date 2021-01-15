package notifications

import (
	"fmt"
	"log"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/utils"
)

// NewMeal sends a push notification about a new meal
func NewMeal(meal *models.Meal, appScheme string) {
	var user models.User
	userQuery := db.Manager.
		Select("name").
		Where("id = ?", meal.UserID).
		Last(&user).
		Error
	if err := userQuery; err != nil {
		log.Println(err)
	}

	title := "Meal Planned"
	body := fmt.Sprintf("%v added a meal to your meal plan", user.Name)
	deviceTokens, err := MealUserTokens(meal)
	if err != nil {
		log.Println(err)
	}
	for i := range deviceTokens {
		Send(title, body, deviceTokens[i], "Meal", meal.ID.String(), appScheme)
	}
}

// MealRemoved sends a push notification about a meal being removed from a meal plan
func MealRemoved(meal models.Meal, appScheme string) {
	var user models.User
	userQuery := db.Manager.
		Select("name").
		Where("id = ?", meal.UserID).
		Last(&user).
		Error
	if err := userQuery; err != nil {
		log.Println(err)
	}

	title := "Meal Removed"
	nameTruncated := utils.TruncateString(meal.Name, 12)
	body := fmt.Sprintf("%v removed \"%v\" from your meal plan", user.Name, nameTruncated)
	deviceTokens, err := MealUserTokens(&meal)
	if err != nil {
		log.Println(err)
	}
	for i := range deviceTokens {
		Send(title, body, deviceTokens[i], "Meal", meal.ID.String(), appScheme)
	}
}

// MealUserTokens fetches apns device tokens for all meal users
func MealUserTokens(meal *models.Meal) (tokens []string, err error) {
	var mealUsers []models.MealUser
	mealUsersQuery := db.Manager.
		Select("meal_users.user_id").
		Joins("INNER JOIN meals ON meals.id = meal_users.meal_id").
		Joins("INNER JOIN stores ON stores.id = meals.store_id").
		Joins("INNER JOIN store_users ON store_users.store_id = stores.id").
		Joins("INNER JOIN store_user_preferences ON store_user_preferences.store_user_id = store_users.id").
		Where("meal_users.meal_id = ?", meal.ID).
		Where("meal_users.user_id NOT IN (meals.user_id)").
		Where("store_user_preferences.notifications = ?", true).
		Group("meal_users.user_id").
		Find(&mealUsers).
		Error
	if err := mealUsersQuery; err != nil {
		return tokens, err
	}

	// Fetch the tokens for each store user
	var t []string
	for i := range mealUsers {
		userTokens, err := DeviceTokensForUser(mealUsers[i].UserID)
		if err != nil {
			return tokens, err
		}
		t = append(t, userTokens...)
	}
	return t, nil
}

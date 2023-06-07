package meals

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/notifications"
	uuid "github.com/satori/go.uuid"
)

// DeleteMeal deletes a meal by ID
func DeleteMeal(mealID uuid.UUID, userID uuid.UUID, appScheme string) (meal models.Meal, err error) {
	query := db.Manager.
		Joins("INNER JOIN meal_users ON meal_users.meal_id = meals.id").
		Where("meals.id = ?", mealID).
		Where("meal_users.user_id = ?", userID).
		Last(&meal).
		Error
	if err := query; err != nil {
		return meal, err
	}

	// Send push notification before deletion so that we still have meal users
	//
	// Note: we send push notification here instead of at the resolver level like
	// other cases because we're deleting; when the resolver returns the meal object,
	// it appears it no longer has meal users associated.
	//
	// FIXME: would be nice to find a way around this and do this in the resolver...
	// feels wrong to do this inside the package
	go notifications.MealRemoved(meal, appScheme)

	if err := db.Manager.Where("id = ?", meal.ID).Delete(&meal).Error; err != nil {
		return meal, err
	}
	return meal, nil
}

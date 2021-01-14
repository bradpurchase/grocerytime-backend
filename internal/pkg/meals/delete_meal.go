package meals

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// DeleteMeal deletes a meal by ID
func DeleteMeal(mealID interface{}, userID uuid.UUID) (deletedMeal models.Meal, err error) {
	var meal models.Meal
	query := db.Manager.
		Where("meals.id = ?", mealID).
		Where("meals.user_id = ?", userID).
		Last(&meal).
		Error
	if err := query; err != nil {
		return deletedMeal, err
	}
	if err := db.Manager.Where("id = ?", meal.ID).Delete(&meal).Error; err != nil {
		return deletedMeal, err
	}
	return meal, nil
}

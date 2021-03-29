package meals

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// DeleteRecipe deletes a recipe
func DeleteRecipe(recipeID interface{}, userID uuid.UUID) (recipe models.Recipe, err error) {
	query := db.Manager.
		Where("id = ? AND user_id = ?", recipeID, userID).
		Last(&recipe).
		Error
	if err := query; err != nil {
		return recipe, err
	}

	// Check if this recipe is in the user's upcoming meal plan and prevent delete
	var upcomingMealsCount int64
	upcomingMealsExistQuery := db.Manager.
		Model(&models.Meal{}).
		Where("recipe_id = ?", recipeID).
		Where("user_id = ?", userID).
		Where("date::date >= current_date").
		Count(&upcomingMealsCount).
		Error
	if err := upcomingMealsExistQuery; err != nil {
		return recipe, err
	}
	if upcomingMealsCount > 0 {
		return recipe, errors.New("cannot delete recipe because there are upcoming meals planned for it")
	}

	if err := db.Manager.Where("id = ? AND user_id = ?", recipeID, userID).Delete(&recipe).Error; err != nil {
		return recipe, err
	}

	return recipe, nil
	//TODO TEST THIS
}

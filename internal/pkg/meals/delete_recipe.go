package meals

import (
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

	if err := db.Manager.Where("id = ? AND user_id = ?", recipeID, userID).Delete(&recipe).Error; err != nil {
		return recipe, err
	}

	return recipe, nil
	//TODO TEST THIS
}

package meals

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// RetrieveRecipes retrieves recipes added by userID
func RetrieveRecipes(userID uuid.UUID) (recipes []models.Recipe, err error) {
	query := db.Manager.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&recipes).
		Error
	if err := query; err != nil {
		return recipes, err
	}
	return recipes, nil
}

// RetrieveRecipe retrieves a recipe by ID
func RetrieveRecipe(id interface{}) (recipe models.Recipe, err error) {
	query := db.Manager.
		Preload("Ingredients").
		Where("id = ?", id).
		Order("created_at DESC").
		Last(&recipe).
		Error
	if err := query; err != nil {
		return recipe, err
	}
	return recipe, nil
}

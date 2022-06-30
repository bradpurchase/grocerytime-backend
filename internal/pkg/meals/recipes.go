package meals

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// RetrieveRecipes retrieves recipes added by userID
func RetrieveRecipes(userID uuid.UUID, args map[string]interface{}) (recipes []models.Recipe, err error) {
	query := db.Manager.Preload("Ingredients").Where("user_id = ?", userID)
	if args["mealType"] != nil && args["mealType"] != "" {
		query = query.Where("meal_type = ?", args["mealType"].(string))
	}
	if args["limit"] != nil {
		query = query.Limit(args["limit"].(int))
	}
	query = query.Order("created_at DESC").Find(&recipes)
	if err := query.Error; err != nil {
		return recipes, err
	}
	return recipes, nil
}

// RetrieveRecipe retrieves a recipe by ID
func RetrieveRecipe(id uuid.UUID) (recipe models.Recipe, err error) {
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

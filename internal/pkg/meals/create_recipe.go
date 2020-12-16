package meals

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// CreateRecipe creates a recipe record and associated  records
func CreateRecipe(userID uuid.UUID, args map[string]interface{}) (recipe *models.Recipe, err error) {
	var url string
	if args["url"] != nil {
		url = args["url"].(string)
	}
	var mealType string
	if args["mealType"] != nil {
		mealType = args["mealType"].(string)
	}
	recipe = &models.Recipe{
		UserID:   userID,
		Name:     args["name"].(string),
		URL:      &url,
		MealType: &mealType,
	}
	if err := db.Manager.Create(&recipe).Error; err != nil {
		return recipe, err
	}

	// Handle ingredients
	recipeID := recipe.ID
	ingredients := args["ingredients"].([]interface{})
	if err := CreateRecipeIngredients(recipeID, ingredients); err != nil {
		return recipe, errors.New("could not create ingredients")
	}

	return recipe, nil
}

// CreateRecipeIngredients creates recipe_ingredients records associated with a RecipeID
func CreateRecipeIngredients(recipeID uuid.UUID, ingredients []interface{}) (err error) {
	for i := range ingredients {
		ingredient := &models.RecipeIngredient{
			RecipeID: recipeID,
			Name:     ingredients[i].(map[string]interface{})["name"].(string),
			Quantity: ingredients[i].(map[string]interface{})["quantity"].(int),
		}
		if err := db.Manager.Create(&ingredient).Error; err != nil {
			return err
		}
	}
	return nil
}

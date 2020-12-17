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

	ingredients, err := CompileRecipeIngredients(args["ingredients"].([]interface{}))
	if err != nil {
		return recipe, errors.New("could not create ingredients")
	}
	recipe = &models.Recipe{
		UserID:      userID,
		Name:        args["name"].(string),
		URL:         &url,
		MealType:    &mealType,
		Ingredients: ingredients,
	}
	if err := db.Manager.Create(&recipe).Error; err != nil {
		return recipe, err
	}
	return recipe, nil
}

// CompileRecipeIngredients compiles []models.RecipeIngredient for insertion in a recipe
func CompileRecipeIngredients(ingArg []interface{}) (ingredients []models.RecipeIngredient, err error) {
	for i := range ingArg {
		amount := ingArg[i].(map[string]interface{})["amount"].(float64)
		unit := ingArg[i].(map[string]interface{})["unit"]
		var unitStr string
		if unit != nil {
			unitStr = unit.(string)
		}
		ingredient := models.RecipeIngredient{
			Name:   ingArg[i].(map[string]interface{})["name"].(string),
			Amount: &amount,
			Unit:   &unitStr,
		}
		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

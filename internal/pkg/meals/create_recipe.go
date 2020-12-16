package meals

import (
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
	// TODO: move this to a function.. and try to clean up?
	ingredients := args["ingredients"].([]interface{})
	for i := range ingredients {
		amount := ingredients[i].(map[string]interface{})["amount"].(int) // TODO: needs to be float (e.g. 3.5 tablespoons)
		unit := ingredients[i].(map[string]interface{})["unit"]
		var unitStr string
		if unit != nil {
			unitStr = unit.(string)
		}
		ingredient := &models.RecipeIngredient{
			RecipeID: recipe.ID,
			Name:     ingredients[i].(map[string]interface{})["name"].(string),
			Amount:   &amount,
			Unit:     &unitStr,
			Quantity: ingredients[i].(map[string]interface{})["quantity"].(int),
		}
		if err := db.Manager.Create(&ingredient).Error; err != nil {
			return recipe, err
		}
	}

	return recipe, nil
}

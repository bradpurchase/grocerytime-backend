package meals

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// CreateRecipe creates a recipe record and associated  records
func CreateRecipe(args map[string]interface{}) (recipe *models.Recipe, err error) {
	var url string
	if args["url"] != nil {
		url = args["url"].(string)
	}

	var mealType string
	if args["mealType"] != nil {
		mealType = args["mealType"].(string)
	}
	recipe = &models.Recipe{
		Name:     args["name"].(string),
		URL:      &url,
		MealType: &mealType,
	}
	if err := db.Manager.Create(&recipe).Error; err != nil {
		return recipe, err
	}

	// Handle ingredients
	//ingredients := args["ingredients"]

	// Handle recipe_users (current user)
}

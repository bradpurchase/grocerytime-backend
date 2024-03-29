package meals

import (
	"encoding/json"
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
)

// CreateRecipe creates a recipe record and associated records
func CreateRecipe(userID uuid.UUID, args map[string]interface{}) (recipe *models.Recipe, err error) {
	// Optional fields
	var mealType string
	if args["mealType"] != nil {
		mealType = args["mealType"].(string)
	}
	var url string
	if args["url"] != nil {
		url = args["url"].(string)
	}
	var imageURL string
	if args["imageUrl"] != nil {
		imageURL = args["imageUrl"].(string)
	}
	var desc string
	if args["description"] != nil {
		desc = args["description"].(string)
	}

	ingredientsArg := args["ingredients"]
	var ingredients []models.RecipeIngredient
	if ingredientsArg != nil {
		ingredients, err = CompileRecipeIngredients(ingredientsArg.([]interface{}))
		if err != nil {
			return recipe, errors.New("could not create ingredients")
		}
	}

	var instructions []interface{}
	if args["instructions"] != nil {
		instructions = args["instructions"].([]interface{})
	}
	instructionsEncoded, err := json.Marshal(instructions)
	if err != nil {
		return recipe, err
	}
	instructionsJSON := datatypes.JSON(instructionsEncoded)

	recipe = &models.Recipe{
		UserID:       userID,
		Name:         args["name"].(string),
		Description:  &desc,
		MealType:     &mealType,
		URL:          &url,
		ImageURL:     &imageURL,
		Ingredients:  ingredients,
		Instructions: &instructionsJSON,
	}
	if err := db.Manager.Create(&recipe).Error; err != nil {
		return recipe, err
	}
	return recipe, nil
}

// CompileRecipeIngredients compiles []models.RecipeIngredient for insertion in a recipe
func CompileRecipeIngredients(ingArg []interface{}) (ingredients []models.RecipeIngredient, err error) {
	for i := range ingArg {
		ing := ingArg[i].(map[string]interface{})

		var amount string
		if ing["amount"] != nil {
			amount = ing["amount"].(string)
		}
		var unit string
		if ing["unit"] != nil {
			unit = ing["unit"].(string)
		}

		var notes string
		if ing["notes"] != nil {
			notes = ing["notes"].(string)
		}
		ingredient := models.RecipeIngredient{
			Name:   ing["name"].(string),
			Amount: &amount,
			Unit:   &unit,
			Notes:  &notes,
		}
		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

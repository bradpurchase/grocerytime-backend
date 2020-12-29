package meals

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// PlanMeal creates a meal record and associated records
func PlanMeal(userID uuid.UUID, args map[string]interface{}) (meal *models.Meal, err error) {
	var mealType string
	if args["mealType"] != nil {
		mealType = args["mealType"].(string)
	}
	var notes string
	if args["notes"] != nil {
		notes = args["notes"].(string)
	}
	date := args["date"].(string)
	recipeID, err := uuid.FromString(args["recipeId"].(string))
	if err != nil {
		return meal, errors.New("recipeId arg not a UUID")
	}

	meal = &models.Meal{
		RecipeID: recipeID,
		UserID:   userID,
		Name:     args["name"].(string),
		MealType: &mealType,
		Servings: args["servings"].(int),
		Notes:    &notes,
		Date:     date,
	}
	if err := db.Manager.Create(&meal).Error; err != nil {
		return meal, err
	}

	// TODO: scan through args["items"] and add each item to args["storeID"]

	return meal, nil
}

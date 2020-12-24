package meals

import (
	"fmt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// PlanMeal creates a meal record and associated records
func PlanMeal(userID uuid.UUID, args map[string]interface{}) (meal *models.Meal, err error) {
	fmt.Println(args)
	// recipeID := args["recipeID"].(uuid.UUID)
	// name := args["name"].(string)
	// var mealType string
	// if args["mealType"] != nil {
	// 	mealType = args["mealType"].(string)
	// }
	// var notes string
	// if args["notes"] != nil {
	// 	notes = args["notes"].(string)
	// }
	// date := args["date"].(time.Time)
	// meal = &models.Meal{
	// 	RecipeID: recipeID,
	// 	Name:     name,
	// 	MealType: &mealType,
	// 	Notes:    &notes,
	// 	Date:     date,
	// }
	return meal, nil
}

package meals

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// UpdateMeal updates a meal with the arguments provided
func UpdateMeal(args map[string]interface{}) (meal models.Meal, err error) {
	if err := db.Manager.Where("id = ?", args["id"]).First(&meal).Error; err != nil {
		return meal, err
	}

	if args["name"] != nil {
		meal.Name = args["name"].(string)
	}
	if args["mealType"] != nil {
		mealType := args["mealType"].(string)
		meal.MealType = &mealType
	}
	if args["servings"] != nil {
		meal.Servings = args["servings"].(int)
	}
	if args["notes"] != nil {
		notes := args["notes"].(string)
		meal.Notes = &notes
	}
	if args["date"] != nil {
		meal.Date = args["date"].(string)
	}

	if err := db.Manager.Save(&meal).Error; err != nil {
		return meal, err
	}
	return meal, nil
}

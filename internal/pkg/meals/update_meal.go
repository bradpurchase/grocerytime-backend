package meals

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// UpdateMeal updates an item by itemID
func UpdateMeal(args map[string]interface{}, appScheme string) (interface{}, error) {
	var meal models.Meal
	if err := db.Manager.Where("id = ?", args["id"]).First(&meal).Error; err != nil {
		return nil, err
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
		return nil, err
	}
	return meal, nil
}

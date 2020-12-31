package meals

import (
	"errors"
	"fmt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/trips"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

	db.Manager.Transaction(func(tx *gorm.DB) error {
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
			return err
		}

		// Add the associated items to the current trip in the store
		storeName := args["storeName"].(string)
		items := args["items"].([]interface{})
		_, e := AddMealIngredientsToStore(storeName, userID, meal.ID, items)
		if e != nil {
			return e
		}

		// TODO: populate meal_users

		return nil
	})

	return meal, nil
}

// AddMealIngredientsToStore will add the items associated with this meal to the user's selected store
func AddMealIngredientsToStore(storeName string, userID uuid.UUID, mealID uuid.UUID, itemsArg []interface{}) (addedItems []*models.Item, err error) {
	var items []interface{}
	for i := range itemsArg {
		item := itemsArg[i].(map[string]interface{})
		quantity := item["quantity"].(int)
		if quantity > 0 {
			items = append(items, fmt.Sprintf("%s x %d", item["name"], item["quantity"]))
		}
	}
	args := map[string]interface{}{
		"storeName": storeName,
		"items":     items,
	}
	itemsAdded, err := trips.AddItemsToStore(userID, args)
	if err != nil {
		return addedItems, err
	}
	return itemsAdded, nil
}

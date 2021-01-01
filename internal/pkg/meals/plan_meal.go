package meals

import (
	"errors"
	"fmt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
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
	storeID, err := uuid.FromString(args["storeId"].(string))
	if err != nil {
		return meal, errors.New("storeId arg not a UUID")
	}

	// Populate meal users by fetching users in the associated store
	storeUsers, err := stores.RetrieveStoreUsers(storeID)
	if err != nil {
		return meal, err
	}
	var mealUsers []models.MealUser
	for _, storeUser := range storeUsers {
		mealUsers = append(mealUsers, models.MealUser{UserID: storeUser.UserID})
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
			Users:    mealUsers,
		}
		if err := db.Manager.Create(&meal).Error; err != nil {
			return err
		}

		// Add the associated items to the current trip in the store
		items := args["items"].([]interface{})
		_, e := AddMealIngredientsToStore(storeID, userID, meal.ID, items)
		if e != nil {
			return e
		}

		return nil
	})

	return meal, nil
}

// AddMealIngredientsToStore will add the items associated with this meal to the user's selected store
func AddMealIngredientsToStore(storeID uuid.UUID, userID uuid.UUID, mealID uuid.UUID, itemsArg []interface{}) (addedItems []*models.Item, err error) {
	var items []interface{}
	for i := range itemsArg {
		item := itemsArg[i].(map[string]interface{})
		quantity := item["quantity"].(int)
		if quantity > 0 {
			// TODO: attribute meal_id somehow (probably need to refactor AddItemsToStore "items" arg to support quantity and meal_id etc)
			items = append(items, fmt.Sprintf("%s x %d", item["name"], item["quantity"]))
		}
	}

	// Fetch store name
	var store models.Store
	if err := db.Manager.Select("name").Where("id = ?", storeID).First(&store).Error; err != nil {
		return addedItems, errors.New("store not found for storeId")
	}

	args := map[string]interface{}{
		"storeName": store.Name,
		"items":     items,
	}
	itemsAdded, err := trips.AddItemsToStore(userID, args)
	if err != nil {
		return addedItems, err
	}

	// Update items to attribute meal_id
	//
	// Note: Ideally, I'd like this to be done at the time of adding the items to avoid
	// a second query to update them to attribute meal_id. The challenge here is that
	// the mutation that calls trips.AddItemsToStore cannot be easily updated to modify
	// the shape of its "items" argument (such as to allow a meal_id or quantity parameter),
	// since we cannot rely on every user updating the app.
	var itemIds []uuid.UUID
	for i := range itemsAdded {
		itemIds = append(itemIds, itemsAdded[i].ID)
	}
	updateQuery := db.Manager.
		Model(&models.Item{}).
		Where("id IN (?)", itemIds).
		Update("meal_id", mealID).
		Error
	if err := updateQuery; err != nil {
		return addedItems, err
	}

	return itemsAdded, nil
}

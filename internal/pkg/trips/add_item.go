package trips

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

//go:embed FoodClassification.json
var foods string

// AddItem adds an item to a trip and handles things like permission checks
func AddItem(userID uuid.UUID, args map[string]interface{}) (addedItem *models.Item, err error) {
	tripID := args["tripId"].(uuid.UUID)
	trip := &models.GroceryTrip{}
	if err := db.Manager.Select("store_id").Where("id = ?", tripID).Last(&trip).Error; err != nil {
		return addedItem, errors.New("trip does not exist")
	}

	// Verify that the current user belongs to this store
	storeUser := &models.StoreUser{}
	if err := db.Manager.Where("store_id = ? AND user_id = ?", trip.StoreID, userID).First(&storeUser).Error; err != nil {
		return addedItem, errors.New("user does not belong to this store")
	}

	itemCompleted := false
	itemName := args["name"].(string)
	quantity := 1
	if args["quantity"] != nil {
		quantity = args["quantity"].(int)
	}

	item := &models.Item{
		GroceryTripID: tripID,
		UserID:        userID,
		Name:          itemName,
		Quantity:      quantity,
		Position:      1,
		Completed:     &itemCompleted,
	}

	if args["stapleItemId"] != nil {
		stapleItemID := args["stapleItemId"].(uuid.UUID)
		item.StapleItemID = &stapleItemID
	}

	// If categoryName is explicitly provided in the arguments, use it,
	// otherwise we need to determine it automagically ✨
	var categoryName string
	if args["categoryName"] != nil {
		categoryName = args["categoryName"].(string)
	} else {
		categoryName, err = DetermineCategoryName(itemName, trip.StoreID)
		if err != nil {
			return item, err
		}
	}

	category, err := FetchGroceryTripCategory(tripID, categoryName)
	if err != nil {
		return addedItem, errors.New("could not find or create grocery trip category")
	}
	item.CategoryID = &category.ID

	if err := db.Manager.Create(&item).Error; err != nil {
		return addedItem, err
	}
	return item, nil
}

//TODO TEST THIS

// FetchGroceryTripCategory retrieves a grocery trip category for a new item by
// finding or creating a category depending on if one exists by the name provided
func FetchGroceryTripCategory(tripID uuid.UUID, name string) (category models.GroceryTripCategory, err error) {
	groceryTripCategory := models.GroceryTripCategory{}
	query := db.Manager.
		Select("grocery_trip_categories.id").
		Joins("INNER JOIN store_categories ON store_categories.id = grocery_trip_categories.store_category_id").
		Where("grocery_trip_categories.grocery_trip_id = ?", tripID).
		Where("store_categories.name = ?", name).
		First(&groceryTripCategory).
		Error
	if err := query; errors.Is(err, gorm.ErrRecordNotFound) {
		newCategory, err := CreateGroceryTripCategory(tripID, name)
		if err != nil {
			return category, err
		}
		return newCategory, err
	}
	return groceryTripCategory, nil
}

// CreateGroceryTripCategory creates a grocery trip category by name
func CreateGroceryTripCategory(tripID uuid.UUID, name string) (category models.GroceryTripCategory, err error) {
	storeCategory := models.StoreCategory{}
	query := db.Manager.
		Select("store_categories.id").
		Joins("INNER JOIN stores ON stores.id = store_categories.store_id").
		Joins("INNER JOIN grocery_trips ON grocery_trips.store_id = stores.id").
		Where("store_categories.name = ?", name).
		Where("grocery_trips.id = ?", tripID).
		First(&storeCategory).
		Error
	if err := query; err != nil {
		return category, errors.New("could not find store category")
	}
	newCategory := models.GroceryTripCategory{
		GroceryTripID:   tripID,
		StoreCategoryID: storeCategory.ID,
	}
	if err := db.Manager.Create(&newCategory).Error; err != nil {
		return category, errors.New("could not create trip category")
	}
	return newCategory, nil
}

// DetermineCategoryName first checks to see if this item's preferred category
// has been saved in the store settings and uses that if so.
// As a fallback, it opens the FoodClassification.json file and scans it
func DetermineCategoryName(name string, storeID uuid.UUID) (result string, err error) {
	result = "Misc."
	name = strings.ToLower(name) // for case-insensitivity

	// Look for the category in store_item_category_settings
	var settings models.StoreItemCategorySettings
	query := db.Manager.
		Where("store_id = ?", storeID).
		Where(datatypes.JSONQuery("items").HasKey(name)).
		First(&settings).
		Error
	if !errors.Is(query, gorm.ErrRecordNotFound) {
		if err := query; err != nil {
			return result, err
		}
		itemSettings := settings.Items
		var settingsMap map[string]interface{}
		if err := json.Unmarshal(itemSettings, &settingsMap); err != nil {
			return result, err
		}
		if settingsMap[name] != nil {
			// There is an assigned storeCategoryID in settings for this item.
			// From this we need to find the name of the category and return it
			storeCategoryID, err := uuid.FromString(settingsMap[name].(string))
			if err != nil {
				return result, err
			}
			return FindStoreCategoryName(storeCategoryID), nil
		}
	}

	// Use gjson to quickly fetch it from the embedded FoodClassification.json file
	properName := strings.TrimSpace(name)
	search := fmt.Sprintf("foods.#(text%%\"%s*\").label", properName)
	value := gjson.Get(foods, search)
	foundCategory := value.String()
	if len(foundCategory) > 0 {
		return foundCategory, nil
	}

	return result, nil
}

func FindStoreCategoryName(id uuid.UUID) (name string) {
	var storeCategory models.StoreCategory
	query := db.Manager.
		Select("store_categories.name").
		Where("store_categories.id = ?", id).
		First(&storeCategory).
		Error
	if err := query; err != nil {
		return "Misc."
	}
	return storeCategory.Name
}

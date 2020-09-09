package trips

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

// AddItem adds an item to a trip and handles things like permission checks
func AddItem(db *gorm.DB, userID uuid.UUID, args map[string]interface{}) (interface{}, error) {
	tripID := args["tripId"]
	trip := &models.GroceryTrip{}
	if err := db.Where("id = ?", tripID).Find(&trip).Error; err != nil {
		return nil, err
	}

	// Verify that the current user belongs to this store
	storeUser := &models.StoreUser{}
	if err := db.Where("store_id = ? AND user_id = ?", trip.StoreID, userID).First(&storeUser).Error; err != nil {
		return nil, err
	}

	itemCompleted := false
	itemName := args["name"].(string)
	quantity := args["quantity"].(int)
	item := models.Item{
		GroceryTripID: trip.ID,
		UserID:        userID,
		Name:          itemName,
		Quantity:      quantity,
		Position:      1,
		Completed:     &itemCompleted,
	}

	// If categoryName is explicily provided in the arguments, use it,
	// otherwise we need to determine it automagically ✨
	var categoryName string
	if args["categoryName"] != nil {
		categoryName = args["categoryName"].(string)
	} else {
		categoryName = DetermineCategoryName(itemName)
	}

	category, err := FetchGroceryTripCategory(db, trip.ID, categoryName)
	if err != nil {
		return nil, err
	}
	item.CategoryID = &category.ID

	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

//TODO TEST THIS

// FetchGroceryTripCategory retrieves a grocery trip category for a new item by
// finding or creating a category depending on if one exists by the name provided
func FetchGroceryTripCategory(db *gorm.DB, tripID uuid.UUID, name string) (models.GroceryTripCategory, error) {
	category := models.GroceryTripCategory{}
	query := db.
		Select("grocery_trip_categories.id").
		Joins("INNER JOIN store_categories ON store_categories.id = grocery_trip_categories.store_category_id").
		Where("grocery_trip_categories.grocery_trip_id = ?", tripID).
		Where("store_categories.name = ?", name).
		First(&category).
		Error
	if err := query; errors.Is(err, gorm.ErrRecordNotFound) {
		newCategory, err := CreateGroceryTripCategory(db, tripID, name)
		if err != nil {
			return models.GroceryTripCategory{}, errors.New("could not find or create grocery trip category")
		}
		return newCategory, err
	}
	return category, nil
}

// CreateGroceryTripCategory creates a grocery trip category by name
func CreateGroceryTripCategory(db *gorm.DB, tripID uuid.UUID, name string) (models.GroceryTripCategory, error) {
	storeCategory := models.StoreCategory{}
	query := db.Select("id").Where("name = ?", name).First(&storeCategory).Error
	if err := query; err != nil {
		return models.GroceryTripCategory{}, err
	}
	newCategory := models.GroceryTripCategory{
		GroceryTripID:   tripID,
		StoreCategoryID: storeCategory.ID,
	}
	if err := db.Create(&newCategory).Error; err != nil {
		return models.GroceryTripCategory{}, err
	}
	return newCategory, nil
}

// DetermineCategoryName opens the FoodClassification.json file and
// scans it for the classification associated with the food item very quickly using gson
func DetermineCategoryName(name string) string {
	properName := strings.ToLower(name)

	absPath, _ := filepath.Abs("internal/pkg/trips/data/FoodClassification.json")
	json, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(json)

	search := fmt.Sprintf("foods.#(text%%\"%s*\").label", properName)
	value := gjson.Get(jsonString, search)
	foundCategory := value.String()
	if len(foundCategory) > 0 {
		return foundCategory
	}
	return "Misc."
}

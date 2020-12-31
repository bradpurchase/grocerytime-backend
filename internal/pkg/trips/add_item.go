package trips

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

// AddItem adds an item to a trip and handles things like permission checks
func AddItem(userID uuid.UUID, args map[string]interface{}) (addedItem *models.Item, err error) {
	tripID := args["tripId"]
	// TODO: also accept trip to avoid this query
	// For example, in AddItemsToStore, we fetch the trip, then fetch it again here by ID. this is pointless
	trip := &models.GroceryTrip{}
	if err := db.Manager.Where("id = ?", tripID).Find(&trip).Error; err != nil {
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
	itemName, quantity = ParseItemName(itemName, quantity)

	item := &models.Item{
		GroceryTripID: trip.ID,
		UserID:        userID,
		Name:          itemName,
		Quantity:      quantity,
		Position:      1,
		Completed:     &itemCompleted,
	}

	// If categoryName is explicitly provided in the arguments, use it,
	// otherwise we need to determine it automagically ✨
	var categoryName string
	if args["categoryName"] != nil {
		categoryName = args["categoryName"].(string)
	} else {
		categoryName = DetermineCategoryName(itemName)
	}

	category, err := FetchGroceryTripCategory(trip.ID, categoryName)
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
	query := db.Manager.Select("id").Where("name = ?", name).First(&storeCategory).Error
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

// DetermineCategoryName opens the FoodClassification.json file and
// scans it for the classification associated with the food item very quickly using gson
func DetermineCategoryName(name string) string {
	jsonFile, err := os.Open("data/FoodClassification.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	json, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(json)

	properName := strings.TrimSpace(strings.ToLower(name))
	search := fmt.Sprintf("foods.#(text%%\"%s*\").label", properName)
	value := gjson.Get(jsonString, search)
	foundCategory := value.String()
	if len(foundCategory) > 0 {
		return foundCategory
	}
	return "Misc."
}

// ParseItemName handles inline quantity in the item name (e.g. Orange x 5) and
// returns a parsed version of both the name and quantity
func ParseItemName(name string, quantity int) (parsedName string, parsedQuantity int) {
	re := regexp.MustCompile("^(.*)(\\s)x(\\s?)(\\d+)(\\s+)?")
	match := re.FindStringSubmatch(name)
	if match != nil {
		var err error
		parsedQuantity, err = strconv.Atoi(match[4])
		if err != nil {
			return name, quantity
		}
		// Strip the quantity out of the name
		parsedName = re.ReplaceAllString(name, "$1")
		return parsedName, parsedQuantity
	}
	return name, quantity
}

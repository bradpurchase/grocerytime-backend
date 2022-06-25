package models

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

//go:embed FoodClassification.json
var foods string

// Item defines the model for items
type Item struct {
	ID            uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	GroceryTripID uuid.UUID  `gorm:"type:uuid;not null;index:idx_items_grocery_trip_id_name"`
	CategoryID    *uuid.UUID `gorm:"type:uuid;not null"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null"`
	StapleItemID  *uuid.UUID `gorm:"type:uuid;index"`
	Name          string     `gorm:"type:varchar(100);not null;index:idx_items_grocery_trip_id_name"`
	Quantity      int        `gorm:"default:1;not null"`
	Completed     *bool      `gorm:"default:false;not null"`
	Position      int        `gorm:"default:1;not null"`
	Notes         *string    `gorm:"type:varchar(255)"`
	MealID        *uuid.UUID `gorm:"type:uuid"`
	MealName      *string    `gorm:"type:varchar(255)"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Associations
	GroceryTrip GroceryTrip
	Meal        Meal
	StapleItem  StoreStapleItem
}

// BeforeCreate hook updates the item position
func (i *Item) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Exec("UPDATE items SET position = position + 1 WHERE grocery_trip_id = ? AND position >= 0", i.GroceryTripID)
	return nil
}

// AfterCreate hook to touch the associated grocery trip after an item is created
// so that its UpdatedAt column is updated
func (i *Item) AfterCreate(tx *gorm.DB) (err error) {
	tx.Model(&GroceryTrip{}).Where("id = ?", i.GroceryTripID).Update("updated_at", time.Now())
	return nil
}

// BeforeSave hook
func (i *Item) BeforeSave(tx *gorm.DB) (err error) {
	// Verify that the item can be added/saved to the trip
	var trip GroceryTrip
	if err := tx.Select("store_id").Where("id = ?", i.GroceryTripID).Last(&trip).Error; err != nil {
		return err
	}
	var storeUser StoreUser
	if err := tx.Where("store_id = ? AND user_id = ?", trip.StoreID, i.UserID).First(&storeUser).Error; err != nil {
		return errors.New("user does not belong to this store")
	}

	// Parse the item name and quantity
	i.Name, i.Quantity = i.parseItemName()

	// Determine the proper category for the item
	categoryName, err := i.DetermineCategoryName(trip.StoreID, tx)
	if err != nil {
		return err
	}
	category, err := i.FetchGroceryTripCategory(categoryName, tx)
	if err != nil {
		return errors.New("could not find or create grocery trip category")
	}
	i.CategoryID = &category.ID

	return nil
}

// BeforeUpdate hook handles reordering items
func (i *Item) BeforeUpdate(tx *gorm.DB) (err error) {
	item := &Item{}
	if err := tx.Where("id = ?", i.ID).Find(&item).Error; err != nil {
		return err
	}
	currPosition := item.Position
	newPosition := i.Position
	if currPosition == newPosition {
		return nil
	}
	if currPosition > newPosition {
		tx.Exec("UPDATE items SET position = position + 1 WHERE grocery_trip_id = ? AND position >= ? AND position < ?", i.GroceryTripID, newPosition, currPosition)
	} else {
		tx.Exec("UPDATE items SET position = position - 1 WHERE grocery_trip_id = ? AND position > ? AND position <= ?", i.GroceryTripID, currPosition, newPosition)
	}
	return nil
}

func (i *Item) AfterUpdate(tx *gorm.DB) (err error) {
	tx.Model(&GroceryTrip{}).Where("id = ?", i.GroceryTripID).Update("updated_at", time.Now())
	return nil
}

// parseItemName handles inline quantity in the item name (e.g. Orange x 5) and
// returns a parsed version of both the name and quantity
func (i *Item) parseItemName() (parsedName string, parsedQuantity int) {
	re := regexp.MustCompile("^(.*)(\\s)x(\\s?)(\\d+)(\\s+)?")
	match := re.FindStringSubmatch(i.Name)
	if match != nil {
		var err error
		parsedQuantity, err = strconv.Atoi(match[4])
		if err != nil {
			return i.Name, i.Quantity
		}
		// Strip the quantity out of the name
		parsedName = re.ReplaceAllString(i.Name, "$1")
		return parsedName, parsedQuantity
	}
	return i.Name, i.Quantity
}

// DetermineCategoryName first checks to see if this item's preferred category
// has been saved in the store settings and uses that if so.
// As a fallback, it opens the FoodClassification.json file and scans it
func (i *Item) DetermineCategoryName(storeID uuid.UUID, tx *gorm.DB) (result string, err error) {
	result = "Misc."
	name := strings.ToLower(i.Name) // for case-insensitivity

	// Look for the category in store_item_category_settings
	var settings StoreItemCategorySettings
	query := tx.
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
			return FindStoreCategoryName(storeCategoryID, tx), nil
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

func FindStoreCategoryName(id uuid.UUID, tx *gorm.DB) (name string) {
	var storeCategory StoreCategory
	query := tx.
		Select("store_categories.name").
		Where("store_categories.id = ?", id).
		First(&storeCategory).
		Error
	if err := query; err != nil {
		return "Misc."
	}
	return storeCategory.Name
}

// FetchGroceryTripCategory retrieves a grocery trip category for a new item by
// finding or creating a category depending on if one exists by the name provided
func (i *Item) FetchGroceryTripCategory(name string, tx *gorm.DB) (category GroceryTripCategory, err error) {
	groceryTripCategory := GroceryTripCategory{}
	query := tx.
		Select("grocery_trip_categories.id").
		Joins("INNER JOIN store_categories ON store_categories.id = grocery_trip_categories.store_category_id").
		Where("grocery_trip_categories.grocery_trip_id = ?", i.GroceryTripID).
		Where("store_categories.name = ?", name).
		First(&groceryTripCategory).
		Error
	if err := query; errors.Is(err, gorm.ErrRecordNotFound) {
		newCategory, err := i.CreateGroceryTripCategory(name, tx)
		if err != nil {
			return category, err
		}
		return newCategory, err
	}
	return groceryTripCategory, nil
}

// CreateGroceryTripCategory creates a grocery trip category by name
func (i *Item) CreateGroceryTripCategory(name string, tx *gorm.DB) (category GroceryTripCategory, err error) {
	storeCategory := StoreCategory{}
	query := tx.
		Select("store_categories.id").
		Joins("INNER JOIN stores ON stores.id = store_categories.store_id").
		Joins("INNER JOIN grocery_trips ON grocery_trips.store_id = stores.id").
		Where("store_categories.name = ?", name).
		Where("grocery_trips.id = ?", i.GroceryTripID).
		First(&storeCategory).
		Error
	if err := query; err != nil {
		return category, errors.New("could not find store category")
	}
	newCategory := GroceryTripCategory{
		GroceryTripID:   i.GroceryTripID,
		StoreCategoryID: storeCategory.ID,
	}
	if err := tx.Create(&newCategory).Error; err != nil {
		return category, errors.New("could not create trip category")
	}
	return newCategory, nil
}

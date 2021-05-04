package trips

import (
	"encoding/json"
	"strings"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
)

// UpdateItem updates an item by itemID
func UpdateItem(args map[string]interface{}) (interface{}, error) {
	item := &models.Item{}
	if err := db.Manager.Where("id = ?", args["itemId"]).First(&item).Error; err != nil {
		return nil, err
	}

	if args["name"] != nil {
		item.Name = args["name"].(string)
	}
	if args["completed"] != nil {
		completed := args["completed"].(bool)
		item.Completed = &completed
		//item.Position = GetNewPosition(db, item.GroceryTripID, completed)
	}
	if args["quantity"] != nil {
		item.Quantity = args["quantity"].(int)
	}
	if args["position"] != nil {
		item.Position = args["position"].(int)
	}
	if args["notes"] != nil {
		notes := args["notes"].(string)
		item.Notes = &notes
	}
	if args["storeCategoryId"] != nil {
		storeCategoryID, _ := uuid.FromString(args["storeCategoryId"].(string))
		groceryTripCategory := models.GroceryTripCategory{
			GroceryTripID:   item.GroceryTripID,
			StoreCategoryID: storeCategoryID,
		}
		// This returns nothing unless the category is already present in the trip,
		// so we do a FirstOrCreate here
		if err := db.Manager.Where(groceryTripCategory).FirstOrCreate(&groceryTripCategory).Error; err != nil {
			return nil, err
		}
		item.CategoryID = &groceryTripCategory.ID

		// If the user opted to save the store category ID, add it to
		// the store_item_category_settings table
		saveCategory := args["saveStoreCategoryId"]
		if saveCategory != nil && saveCategory.(bool) {
			if err := SaveStoreCategorySelection(item, storeCategoryID); err != nil {
				return nil, err
			}
		}
	}
	if err := db.Manager.Save(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// GetNewPosition gets the new position of an updated item
func GetNewPosition(tripID uuid.UUID, completed bool) int {
	newPosition := 1
	if completed {
		// If the item was marked completed, move to the bottom of the store
		// The BeforeUpdate hook on items will handle reordering the items around it
		bottomItem := &models.Item{}
		db.Manager.
			Select("position").
			Where("grocery_trip_id = ?", tripID).
			Order("position DESC").
			Limit(1).
			Find(&bottomItem)
		newPosition = bottomItem.Position
	}
	return newPosition
}

// SaveStoreCategorySelection updates the store item category settings
// for a given item in a store, so that the item will be added to this category going forwards
func SaveStoreCategorySelection(item *models.Item, storeCategoryID uuid.UUID) (err error) {
	var store models.Store
	query := db.Manager.
		Model(&models.Item{}).
		Select("stores.id").
		Joins("INNER JOIN grocery_trips ON grocery_trips.id = items.grocery_trip_id").
		Joins("INNER JOIN stores ON stores.id = grocery_trips.store_id").
		Where("items.id = ?", item.ID).
		Last(&store).
		Error
	if err := query; err != nil {
		return err
	}

	storeItemCategorySetting := models.StoreItemCategorySettings{StoreID: store.ID}
	settingQuery := db.Manager.Where(storeItemCategorySetting).FirstOrCreate(&storeItemCategorySetting).Error
	if err := settingQuery; err != nil {
		return err
	}

	// Set the value in the store_item_category_settings.items jsonb
	itemName := item.Name
	if err := UpdateStoreItemCategorySettings(itemName, storeCategoryID, storeItemCategorySetting); err != nil {
		return err
	}

	return
}

// UpdateStoreItemCategorySettings unpacks the item settings, adds/updates the key
// by item name, repacks the updated item settings into JSON, and saves it
func UpdateStoreItemCategorySettings(
	itemName string,
	storeCategoryID uuid.UUID,
	storeItemCategorySetting models.StoreItemCategorySettings,
) (err error) {
	settings, err := CompileItemSettingsMap(storeItemCategorySetting.Items)
	if err != nil {
		return err
	}

	itemName = strings.ToLower(itemName)
	settings[itemName] = storeCategoryID
	newItemSettings, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	updateQuery := db.Manager.
		Model(&models.StoreItemCategorySettings{}).
		Where("id = ?", storeItemCategorySetting.ID).
		Update("items", newItemSettings).
		Error
	if err := updateQuery; err != nil {
		return err
	}
	return
}

// CompileItemSettingsMap unmarshals the existing item settings json map,
// or makes a new map if one does not exist yet
func CompileItemSettingsMap(itemSettings datatypes.JSON) (settings map[string]interface{}, err error) {
	if itemSettings == nil {
		settings = make(map[string]interface{})
	}
	if err := json.Unmarshal(itemSettings, &settings); err != nil {
		return nil, err
	}
	return settings, nil
}

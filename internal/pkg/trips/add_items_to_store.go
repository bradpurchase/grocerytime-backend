package trips

import (
	"errors"
	"strings"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/stores"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// AddItemsToStore adds an array of items to a store for a user. It creates
// the store for the user if it doesn't already exist.
func AddItemsToStore(userID uuid.UUID, args map[string]interface{}) (addedItems []*models.Item, err error) {
	var store models.Store
	storeName, val := args["storeName"]
	if val && args["storeName"] != nil {
		store, err = FindOrCreateStore(userID, storeName.(string))
		if err != nil {
			return addedItems, errors.New("could not find or create store")
		}
	} else {
		// TODO: what if there's no default store set? handle this case...
		// could fall back to the user's first store added as a hail mary
		store, err = FindDefaultStore(userID)
		if err != nil {
			return nil, errors.New("could not retrieve default store")
		}
	}

	// Fetch the current trip for this store
	trip, err := RetrieveCurrentStoreTrip(store.ID)
	if err != nil {
		return addedItems, errors.New("could not find current trip in store")
	}

	var errorStrings []string
	itemNames := args["items"].([]interface{})
	for i := range itemNames {
		itemName := itemNames[i].(string)
		item, err := AddItem(userID, map[string]interface{}{
			"tripId":   trip.ID,
			"name":     itemName,
			"quantity": 1,
		})
		if err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
		addedItems = append(addedItems, item)
	}
	if len(errorStrings) > 0 {
		return addedItems, errors.New(strings.Join(errorStrings, "\n"))
	}

	return addedItems, nil
}

// FindOrCreateStore finds or creates a store for a userID by name
func FindOrCreateStore(userID uuid.UUID, name string) (storeRecord models.Store, err error) {
	store := models.Store{}
	storeQuery := db.Manager.
		Select("stores.id").
		Joins("INNER JOIN store_users ON store_users.store_id = stores.id").
		Where("store_users.user_id = ?", userID).
		Where("stores.name = ?", name).
		First(&store).
		Error
	if err := storeQuery; errors.Is(err, gorm.ErrRecordNotFound) {
		newStore, err := stores.CreateStore(userID, name)
		if err != nil {
			return storeRecord, errors.New("could not find or create store")
		}
		return newStore, nil
	}
	return store, nil
}

// FindDefaultStore retrieves the ID of the store that is set as the default for the userID provided
func FindDefaultStore(userID uuid.UUID) (store models.Store, err error) {
	query := db.Manager.
		Select("stores.id").
		Joins("INNER JOIN store_users ON store_users.store_id = stores.id").
		Joins("INNER JOIN store_user_preferences ON store_user_preferences.store_user_id = store_users.id").
		Where("store_users.user_id = ?", userID).
		Where("store_user_preferences.default_store = ?", true).
		Last(&store).
		Error
	if err := query; err != nil {
		return store, err
	}
	return store, nil
}

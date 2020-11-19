package trips

import (
	"errors"
	"strings"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// AddItemsToStore adds an array of items to a store for a user. It creates
// the store for the user if it doesn't already exist.
func AddItemsToStore(userID uuid.UUID, args map[string]interface{}) (addedItems []*models.Item, err error) {
	var store models.Store
	if args["storeName"] != nil {
		storeName := args["storeName"].(string)
		store, err = FindOrCreateStore(userID, storeName)
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
	tripID, err := FindCurrentTripIDInStore(store.ID)
	if err != nil {
		return addedItems, errors.New("could not find current trip in store")
	}

	var errorStrings []string
	itemNames := args["items"].([]interface{})
	for i := range itemNames {
		itemName := itemNames[i].(string)
		item, err := AddItem(userID, map[string]interface{}{
			"tripId":   tripID,
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
		Where(models.Store{UserID: userID, Name: name}).
		FirstOrCreate(&store).
		Error
	if err := storeQuery; err != nil {
		return storeRecord, errors.New("could not find or create store")
	}
	return store, nil
}

// FindCurrentTripIDInStore retrieves the ID of the most recent trip in the store that hasn't been completed
//
// TODO: DRY this up with trips.RetrieveCurrentStoreTrip
func FindCurrentTripIDInStore(storeID uuid.UUID) (tripID uuid.UUID, err error) {
	trip := models.GroceryTrip{}
	tripQuery := db.Manager.
		Select("id").
		Where("store_id = ? AND completed = ?", storeID, false).
		Order("created_at DESC").
		Last(&trip).
		Error
	if err := tripQuery; err != nil {
		return tripID, err
	}
	return trip.ID, nil
}

// FindDefaultStore retrieves the ID of the store that is set as the default for the userID provided
func FindDefaultStore(userID uuid.UUID) (store models.Store, err error) {
	query := db.Manager.
		Select("stores.id").
		Joins("INNER JOIN store_users ON store_users.store_id = stores.id").
		Joins("INNER JOIN store_user_preferences ON store_user_preferences.store_user_id = store_users.id").
		Where("store_users.user_id = ?", userID).
		Where("store_user_preferences.default_store = ?", true).
		Find(&store).
		Error
	if err := query; err != nil {
		return store, err
	}
	return store, nil
}

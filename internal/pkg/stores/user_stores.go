package stores

import (
	"errors"
	"fmt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// RetrieveUserStores retrieves stores that the userID has created or has been added to
func RetrieveUserStores(user models.User) ([]models.Store, error) {
	stores := []models.Store{}
	query := db.Manager.
		Select("stores.*").
		Joins("INNER JOIN store_users ON store_users.store_id = stores.id").
		Joins("LEFT OUTER JOIN grocery_trips ON grocery_trips.store_id = stores.id").
		Where("store_users.user_id = ?", user.ID).
		Group("stores.id").
		Order("MAX(grocery_trips.updated_at) DESC").
		Find(&stores).
		Error
	if err := query; err != nil {
		return nil, err
	}
	return stores, nil
}

// RetrieveInvitedUserStores retrieves stores that the userID has created or has been added to
func RetrieveInvitedUserStores(user models.User) ([]models.Store, error) {
	stores := []models.Store{}
	query := db.Manager.
		Select("stores.*").
		Joins("INNER JOIN store_users ON store_users.store_id = stores.id").
		Joins("LEFT OUTER JOIN grocery_trips ON grocery_trips.store_id = stores.id").
		Where("store_users.deleted_at IS NULL").
		Where("store_users.email = ?", user.Email).
		Where("store_users.active = ?", false).
		Group("stores.id").
		Order("MAX(grocery_trips.updated_at) DESC").
		Find(&stores).
		Error
	if err := query; err != nil {
		return nil, err
	}
	return stores, nil
}

// RetrieveStoreForUser retrieves a specific store by storeID and userID
func RetrieveStoreForUser(storeID interface{}, userID uuid.UUID) (models.Store, error) {
	store := models.Store{}
	if err := db.Manager.Where("id = ?", storeID).First(&store).Error; err != nil {
		return store, err
	}
	// Check that the passed userID is a member of this store
	storeUser := &models.StoreUser{}
	if err := db.Manager.Where("store_id = ? AND user_id = ?", storeID, userID).First(&storeUser).Error; err != nil {
		return store, err
	}
	return store, nil
}

// RetrieveStoreForUserByName retrieves a specific store by name and userID.
// It is used in resolvers.CreateStoreResolver to determine whether a store with
// a given name already exists in the user's account, to avoid duplicates.
func RetrieveStoreForUserByName(name string, userID uuid.UUID) (models.Store, error) {
	store := models.Store{}
	if err := db.Manager.Where("name = ? AND user_id = ?", name, userID).First(&store).Error; err != nil {
		return store, err
	}
	return store, nil
}

// DeleteStore deletes a store, its associated trips, items, store users,
// and finally notifies the store users that the store has been deleted
//
// Note: this really performs a soft delete for stores and associated models
func DeleteStore(storeID interface{}, userID uuid.UUID) (models.Store, error) {
	store := models.Store{}
	if err := db.Manager.Where("id = ? AND user_id = ?", storeID, userID).First(&store).Error; err != nil {
		return store, errors.New("couldn't retrieve store")
	}
	if err := db.Manager.Delete(&store).Error; err != nil {
		return store, errors.New("couldn't delete store")
	}

	// Delete items in each trip in this store, and then delete the trips themselves
	trips := []models.GroceryTrip{}
	if err := db.Manager.Where("store_id = ?", storeID).Find(&trips).Error; err != nil {
		return store, errors.New("couldn't find trips in store")
	}
	for i := range trips {
		tripID := trips[i].ID

		// Note: we can just use `.Delete` directly here because
		// we don't need to do anything with the items after deletion.
		// for store users we need to fetch, notify, and *then* delete
		items := []models.Item{}
		if err := db.Manager.Where("grocery_trip_id = ?", tripID).Delete(&items).Error; err != nil {
			return store, errors.New("couldn't delete items in this store's trips")
		}

		trip := models.GroceryTrip{}
		if err := db.Manager.Where("id = ?", tripID).Delete(&trip).Error; err != nil {
			return store, fmt.Errorf("couldn't delete trip: %s", tripID)
		}
	}

	storeUsers := []models.StoreUser{}
	if err := db.Manager.Where("store_id = ?", storeID).Find(&storeUsers).Error; err != nil {
		return store, errors.New("couldn't retrieve store users")
	}

	//TODO notify store users that store was deleted (except creator)
	if err := db.Manager.Where("store_id = ?", storeID).Delete(&storeUsers).Error; err != nil {
		return store, errors.New("couldn't delete store users")
	}

	return store, nil
}

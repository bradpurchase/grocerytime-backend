package stores

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// RetrieveUserStores retrieves stores that the userID has created or has been added to
func RetrieveUserStores(user models.User) (stores []models.Store, err error) {
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
func RetrieveInvitedUserStores(user models.User) (stores []models.Store, err error) {
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
		return stores, err
	}
	return stores, nil
}

// RetrieveStoreForUser retrieves a specific store by storeID and userID
func RetrieveStoreForUser(storeID interface{}, userID uuid.UUID) (store models.Store, err error) {
	if err := db.Manager.Where("id = ?", storeID).First(&store).Error; err != nil {
		return store, err
	}
	// Check that the passed userID is a member of this store
	var storeUser models.StoreUser
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

// DeleteStore handles deletion of a store record
// Note: Associated trips, items, store users etc. are deleted in the AfterDelete hook on the model
func DeleteStore(storeID interface{}, userID uuid.UUID) (deletedStore models.Store, err error) {
	var store models.Store
	if err := db.Manager.Where("id = ? AND user_id = ?", storeID, userID).First(&store).Error; err != nil {
		return deletedStore, errors.New("couldn't retrieve store")
	}
	if err := db.Manager.Delete(&store).Error; err != nil {
		return deletedStore, err
	}
	return store, nil
}

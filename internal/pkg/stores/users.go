package stores

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	uuid "github.com/satori/go.uuid"
)

// InviteToStoreByEmail creates a store_users record for this store ID and email
//
// The store user will be considered pending until the invitation is accepted
// by the user in the app, at which point they are associated by userID instead.
func InviteToStoreByEmail(storeID interface{}, invitedEmail string) (models.StoreUser, error) {
	store := &models.Store{}
	if err := db.Manager.Where("id = ?", storeID).First(&store).Error; err != nil {
		return models.StoreUser{}, err
	}

	storeUserActive := false
	storeUser := models.StoreUser{
		StoreID: store.ID,
		Email:   invitedEmail,
		Active:  &storeUserActive,
	}
	if err := db.Manager.Where(storeUser).FirstOrCreate(&storeUser).Error; err != nil {
		return models.StoreUser{}, err
	}
	return storeUser, nil
}

// AddUserToStore properly associates a user with a store by userID by removing
// the email value and adding the userID value
func AddUserToStore(user models.User, storeID interface{}) (newStoreUser *models.StoreUser, err error) {
	store := &models.Store{}
	if err := db.Manager.Where("id = ?", storeID).First(&store).Error; err != nil {
		return newStoreUser, err
	}

	// Get a count of this user's stores to determine if this will be their default
	var numStores int64
	storesCountQuery := db.Manager.
		Model(&models.StoreUser{}).
		Where("user_id = ? AND active = ?", user.ID, true).
		Count(&numStores).
		Error
	if err := storesCountQuery; err != nil {
		return newStoreUser, err
	}
	defaultStore := numStores == 0

	storeUser := &models.StoreUser{}
	updateStoreUserQuery := db.Manager.
		Where("store_id = ? AND email = ?", storeID, user.Email).
		Find(&storeUser).
		Error
	if err := updateStoreUserQuery; err != nil {
		return newStoreUser, err
	}
	storeUser.Email = ""
	storeUser.UserID = user.ID
	storeUserActive := true
	storeUser.Active = &storeUserActive
	storeUser.DefaultStore = defaultStore
	if err := db.Manager.Save(&storeUser).Error; err != nil {
		return newStoreUser, err
	}
	return storeUser, nil
}

// RemoveUserFromStore removes a user from a store either by userID or email, whichever is present
//
// Used for declining a store invite, and simply removing a user from a store
func RemoveUserFromStore(user models.User, storeID interface{}) (interface{}, error) {
	store := &models.Store{}
	if err := db.Manager.Where("id = ?", storeID).First(&store).Error; err != nil {
		return nil, errors.New("store not found")
	}

	storeUser := &models.StoreUser{}
	query := db.Manager.
		Where("store_id = ?", storeID).
		Where("user_id = ?", user.ID).
		Or("email = ?", user.Email).
		Find(&storeUser).
		Error
	if err := query; err != nil {
		return nil, errors.New("store user not found")
	}

	if err := db.Manager.Where("id = ?", &storeUser.ID).Delete(&storeUser).Error; err != nil {
		return nil, err
	}

	// Get the store creator StoreUser record and User record
	creatorStoreUser := &models.StoreUser{}
	if err := db.Manager.Select("user_id").Where("store_id = ? AND creator = ?", storeID, true).First(&creatorStoreUser).Error; err != nil {
		return nil, err
	}
	creatorUser := &models.User{}
	if err := db.Manager.Select("email").Where("id = ?", creatorStoreUser.UserID).First(&creatorUser).Error; err != nil {
		return nil, err
	}

	// If email is present on the StoreUser record, it means this is a pending store invite
	pending := len(storeUser.Email) > 0
	if pending {
		_, err := mailer.SendStoreInviteDeclinedEmail(store.Name, storeUser.Email, creatorUser.Email)
		if err != nil {
			return nil, err
		}
	} else {
		storeUserUser := &models.User{}
		if err := db.Manager.Select("name").Where("id = ?", storeUser.UserID).Find(&storeUserUser).Error; err != nil {
			return nil, err
		}
		_, err := mailer.SendUserLeftStoreEmail(store.Name, storeUserUser.Name, creatorUser.Email)
		if err != nil {
			return nil, err
		}
	}

	return storeUser, nil
}

// RetrieveStoreUsers finds all store users in a store by storeID5
func RetrieveStoreUsers(storeID uuid.UUID) (interface{}, error) {
	var storeUsers []models.StoreUser
	if err := db.Manager.Where("store_id = ?", storeID).Order("created_at ASC").Find(&storeUsers).Error; err != nil {
		return nil, err
	}
	return storeUsers, nil
}

// RetrieveStoreCreator returns the store user who created a given store
func RetrieveStoreCreator(storeID uuid.UUID) (interface{}, error) {
	storeUser := &models.StoreUser{}
	query := db.Manager.
		Where("store_id = ?", storeID).
		Where("creator = ?", true).
		First(&storeUser).
		Error
	if err := query; err != nil {
		return nil, err
	}

	user := &models.User{}
	if err := db.Manager.Where("id = ?", storeUser.UserID).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

package grocerylist

import (
	"errors"
	"fmt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// RetrieveUserLists retrieves lists that the userID has created or has been added to
func RetrieveUserLists(db *gorm.DB, user models.User) ([]models.List, error) {
	lists := []models.List{}
	query := db.
		Select("lists.*").
		Joins("INNER JOIN list_users ON list_users.list_id = lists.id").
		Joins("LEFT OUTER JOIN grocery_trips ON grocery_trips.list_id = lists.id").
		Where("list_users.user_id = ?", user.ID).
		Or("list_users.email = ?", user.Email).
		Group("lists.id").
		Order("MAX(grocery_trips.updated_at) DESC").
		Find(&lists).
		Error
	if err := query; err != nil {
		return nil, err
	}
	return lists, nil
}

// RetrieveListForUser retrieves a specific list by listID and userID
func RetrieveListForUser(db *gorm.DB, listID interface{}, userID uuid.UUID) (models.List, error) {
	list := models.List{}
	if err := db.Where("id = ?", listID).First(&list).Error; err != nil {
		return list, err
	}
	// Check that the passed userID is a member of this list
	listUser := &models.ListUser{}
	if err := db.Where("list_id = ? AND user_id = ?", listID, userID).First(&listUser).Error; err != nil {
		return list, err
	}
	return list, nil
}

// RetrieveListForUserByName retrieves a specific list by name and userID.
// It is used in resolvers.CreateListResolver to determine whether a list with
// a given name already exists in the user's account, to avoid duplicates.
func RetrieveListForUserByName(db *gorm.DB, name string, userID uuid.UUID) (models.List, error) {
	list := models.List{}
	if err := db.Where("name = ? AND user_id = ?", name, userID).First(&list).Error; err != nil {
		return list, err
	}
	return list, nil
}

// RetrieveSharableList retrieves only publicly information about a list.
// It is used on the share web app to display basic info about a list that someone has been invited to
func RetrieveSharableList(db *gorm.DB, listID interface{}) (models.List, error) {
	list := models.List{}
	query := db.
		Select("lists.id, lists.name, lists.user_id").
		Where("lists.id = ?", listID).
		Find(&list).
		Error
	if err := query; err != nil {
		return list, errors.New("list not found")
	}
	return list, nil
}

// DeleteList deletes a list, its associated trips, items, list users,
// and finally notifies the list users that the list has been deleted
//
// Note: this really performs a soft delete for lists and associated models
func DeleteList(db *gorm.DB, listID interface{}, userID uuid.UUID) (models.List, error) {
	list := models.List{}
	if err := db.Where("id = ? AND user_id = ?", listID, userID).First(&list).Error; err != nil {
		return list, errors.New("couldn't retrieve list")
	}
	if err := db.Delete(&list).Error; err != nil {
		return list, errors.New("couldn't delete list")
	}

	// Delete items in each trip in this list, and then delete the trips themselves
	trips := []models.GroceryTrip{}
	if err := db.Where("list_id = ?", listID).Find(&trips).Error; err != nil {
		return list, errors.New("couldn't find trips in list")
	}
	for i := range trips {
		tripID := trips[i].ID

		// Note: we can just use `.Delete` directly here because
		// we don't need to do anything with the items after deletion.
		// for list users we need to fetch, notify, and *then* delete
		items := []models.Item{}
		if err := db.Where("grocery_trip_id = ?", tripID).Delete(&items).Error; err != nil {
			return list, errors.New("couldn't delete items in this list's trips")
		}

		trip := models.GroceryTrip{}
		if err := db.Where("id = ?", tripID).Delete(&trip).Error; err != nil {
			return list, fmt.Errorf("couldn't delete trip: %s", tripID)
		}
	}

	listUsers := []models.ListUser{}
	if err := db.Where("list_id = ?", listID).Find(&listUsers).Error; err != nil {
		return list, errors.New("couldn't retrieve list users")
	}

	//TODO notify list users that list was deleted (except creator)
	if err := db.Where("list_id = ?", listID).Delete(&listUsers).Error; err != nil {
		return list, errors.New("couldn't delete list users")
	}

	return list, nil
}

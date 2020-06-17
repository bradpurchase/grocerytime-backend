package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// DEPRECATED
// AddUserToList adds a user to a list by user ID
//
// If there is already a list_users record for this list and user,
// this function will simply return the record
func AddUserToList(db *gorm.DB, userID uuid.UUID, listID interface{}) (interface{}, error) {
	list := &models.List{}
	if err := db.Where("id = ?", listID).First(&list).Error; err != nil {
		return nil, err
	}

	listUser := models.ListUser{UserID: userID, ListID: list.ID}
	if err := db.Where(listUser).FirstOrCreate(&listUser).Error; err != nil {
		return nil, err
	}
	return listUser, nil
}

// InviteToListByEmail creates a list_users record for this list ID and email
//
// The list user will be considered pending until the invitation is accepted
// by the user in the app, at which point they are associated by userID instead.
func InviteToListByEmail(db *gorm.DB, listID interface{}, email string) (models.ListUser, error) {
	list := &models.List{}
	if err := db.Where("id = ?", listID).First(&list).Error; err != nil {
		return models.ListUser{}, err
	}

	listUser := models.ListUser{ListID: list.ID, Email: email}
	if err := db.Where(listUser).FirstOrCreate(&listUser).Error; err != nil {
		return models.ListUser{}, err
	}
	return listUser, nil
}

// RetrieveListUsers finds all list users in a list by listID
func RetrieveListUsers(db *gorm.DB, listID uuid.UUID) (interface{}, error) {
	listUsers := []models.ListUser{}
	if err := db.Where("list_id = ?", listID).Order("created_at ASC").Find(&listUsers).Error; err != nil {
		return nil, err
	}
	return listUsers, nil
}

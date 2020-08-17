package grocerylist

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/mailer"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// InviteToListByEmail creates a list_users record for this list ID and email
//
// The list user will be considered pending until the invitation is accepted
// by the user in the app, at which point they are associated by userID instead.
func InviteToListByEmail(db *gorm.DB, listID interface{}, invitedEmail string) (models.ListUser, error) {
	list := &models.List{}
	if err := db.Where("id = ?", listID).First(&list).Error; err != nil {
		return models.ListUser{}, err
	}

	listUserActive := false
	listUser := models.ListUser{
		ListID: list.ID,
		Email:  invitedEmail,
		Active: &listUserActive,
	}
	if err := db.Where(listUser).FirstOrCreate(&listUser).Error; err != nil {
		return models.ListUser{}, err
	}
	return listUser, nil
}

// AddUserToList properly associates a user with a list by userID by removing
// the email value and adding the userID value
func AddUserToList(db *gorm.DB, user models.User, listID interface{}) (interface{}, error) {
	list := &models.List{}
	if err := db.Where("id = ?", listID).First(&list).Error; err != nil {
		return nil, err
	}

	listUser := &models.ListUser{}
	updateListUserQuery := db.
		Where("list_id = ? AND email = ?", listID, user.Email).
		Find(&listUser).
		Error
	if err := updateListUserQuery; err != nil {
		return nil, err
	}
	listUser.Email = ""
	listUser.UserID = user.ID
	listUserActive := true
	listUser.Active = &listUserActive
	if err := db.Save(&listUser).Error; err != nil {
		return nil, err
	}
	return listUser, nil
}

// RemoveUserFromList removes a user from a list either by userID or email, whichever is present
//
// Used for declining a list invite, and simply removing a user from a list
func RemoveUserFromList(db *gorm.DB, user models.User, listID interface{}) (interface{}, error) {
	list := &models.List{}
	if err := db.Where("id = ?", listID).First(&list).Error; err != nil {
		return nil, errors.New("list not found")
	}

	listUser := &models.ListUser{}
	query := db.
		Where("list_id = ?", listID).
		Where("user_id = ?", user.ID).
		Or("email = ?", user.Email).
		Find(&listUser).
		Error
	if err := query; err != nil {
		return nil, errors.New("list user not found")
	}

	if err := db.Delete(&listUser).Error; err != nil {
		return nil, err
	}

	// Get the list creator ListUser record and User record
	creatorListUser := &models.ListUser{}
	if err := db.Select("user_id").Where("list_id = ? AND creator = ?", listID, true).First(&creatorListUser).Error; err != nil {
		return nil, err
	}
	creatorUser := &models.User{}
	if err := db.Select("email").Where("id = ?", creatorListUser.UserID).First(&creatorUser).Error; err != nil {
		return nil, err
	}

	// If email is present on the ListUser record, it means this is a pending list invite
	pending := len(listUser.Email) > 0
	if pending {
		_, err := mailer.SendListInviteDeclinedEmail(list.Name, listUser.Email, creatorUser.Email)
		if err != nil {
			return nil, err
		}
	} else {
		listUserUser := &models.User{}
		if err := db.Select("email").Where("id = ?", listUser.UserID).Find(&listUserUser).Error; err != nil {
			return nil, err
		}
		_, err := mailer.SendUserLeftListEmail(list.Name, listUserUser.Email, creatorUser.Email)
		if err != nil {
			return nil, err
		}
	}

	return listUser, nil
}

// RetrieveListUsers finds all list users in a list by listID5
func RetrieveListUsers(db *gorm.DB, listID uuid.UUID) (interface{}, error) {
	listUsers := []models.ListUser{}
	if err := db.Where("list_id = ?", listID).Order("created_at ASC").Find(&listUsers).Error; err != nil {
		return nil, err
	}
	return listUsers, nil
}

// RetrieveListCreator returns the list user who created a given list
func RetrieveListCreator(db *gorm.DB, listID uuid.UUID) (interface{}, error) {
	listUser := &models.ListUser{}
	query := db.
		Where("list_id = ?", listID).
		Where("creator = ?", true).
		First(&listUser).
		Error
	if err := query; err != nil {
		return nil, err
	}

	user := &models.User{}
	if err := db.Where("id = ?", listUser.UserID).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

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

	// If email is present on the ListUser record, it means this is a pending list invite
	pending := len(listUser.Email) > 0
	if pending {
		creatorListUser := &models.ListUser{}
		if err := db.Select("user_id").Where("creator = ?", true).Find(&creatorListUser).Error; err != nil {
			return nil, err
		}
		user := &models.User{}
		if err := db.Where("id = ?", creatorListUser.UserID).Find(&user).Error; err != nil {
			return nil, err
		}
		_, err := mailer.SendListInviteDeclinedEmail(list.Name, listUser.Email, user.Email)
		if err != nil {
			return nil, err
		}
	} else {
		//TODO mailer sent to user removed from list
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

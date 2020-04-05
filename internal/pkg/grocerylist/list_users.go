package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
)

// AddUserToList adds a user to a list by email. It also handles creating
// a new user record if the user doesn't have an account yet.
func AddUserToList(db *gorm.DB, email string, list *models.List) (interface{}, error) {
	user := &models.User{}
	if err := db.Where("email = ?", email).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		// TODO User doesn't exist flow:
		// - Add the user to the list by email
		// - Send an email to them asking them to sign up
		// - When a user signs up and they are in a UserList
		// associated by Email and not UserID, link the UserID to the UserList
		// so that their "membership" to that list is complete
		listUser := &models.ListUser{}
		return listUser, nil
	}

	listUser := models.ListUser{
		UserID: user.ID,
		ListID: list.ID,
	}
	if err := db.Where(listUser).FirstOrCreate(&listUser).Error; err != nil {
		return nil, err
	}
	return listUser, nil
}

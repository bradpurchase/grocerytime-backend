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
		// TODO User doesn't exist; need to create 'em
		listUser := &models.ListUser{}
		return listUser, nil
	}

	listUser := &models.ListUser{
		UserID: user.ID,
		ListID: list.ID,
	}
	if err := db.Create(&listUser).Error; err != nil {
		return nil, err
	}
	return listUser, nil
}

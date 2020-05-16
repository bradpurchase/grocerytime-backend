package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// AddUserToList adds a user to a list by user ID
func AddUserToList(db *gorm.DB, userID string, list *models.List) (interface{}, error) {
	user := &models.User{}
	if err := db.Where("id = ?", userID).First(&user).Error; gorm.IsRecordNotFoundError(err) {
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

// RetrieveListUsers finds all list users in a list by listID
func RetrieveListUsers(db *gorm.DB, listID uuid.UUID) (interface{}, error) {
	listUsers := []models.ListUser{}
	if err := db.Where("list_id = ?", listID).Order("created_at ASC").Find(&listUsers).Error; err != nil {
		return nil, err
	}
	return listUsers, nil
}

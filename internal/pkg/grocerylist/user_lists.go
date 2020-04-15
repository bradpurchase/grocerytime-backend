package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// RetrieveUserLists retrieves lists that the userID has created or has been added to
func RetrieveUserLists(db *gorm.DB, userID uuid.UUID) (interface{}, error) {
	lists := []models.List{}
	query := db.
		Select("lists.*").
		Joins("INNER JOIN list_users ON list_users.list_id = lists.id").
		Where("list_users.user_id = ?", userID).
		Find(&lists).
		Error
	if err := query; err != nil {
		return nil, err
	}
	return lists, nil
}

// RetrieveListForUser retrieves a specific list for a user
func RetrieveListForUser(db *gorm.DB, listID interface{}, userID uuid.UUID) (interface{}, error) {
	list := &models.List{}
	if err := db.Where("id = ? AND user_id = ?", listID, userID).First(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
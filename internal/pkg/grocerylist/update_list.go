package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// UpdateListForUser updates a list for the given userID with the provided args
func UpdateListForUser(db *gorm.DB, userID uuid.UUID, args map[string]interface{}) (interface{}, error) {
	list := &models.List{}
	if err := db.Where("id = ? AND user_id = ?", args["listId"], userID).First(&list).Error; err != nil {
		return nil, err
	}

	if args["name"] != nil {
		list.Name = args["name"].(string)
	}
	if err := db.Save(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

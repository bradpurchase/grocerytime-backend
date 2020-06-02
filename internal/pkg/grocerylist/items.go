package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// RetrieveItemsInList finds all items in a list by listID
func RetrieveItemsInList(db *gorm.DB, listID uuid.UUID) (interface{}, error) {
	items := []models.Item{}
	query := db.
		Where("list_id = ?", listID).
		Order("position ASC").
		Find(&items).
		Error
	if err := query; err != nil {
		return nil, err
	}
	return items, nil
}

// AddItemToList adds an item to a list and handles things like permission checks
// and handles how it should be sorted etc
func AddItemToList(db *gorm.DB, userID uuid.UUID, args map[string]interface{}) (interface{}, error) {
	listID := args["listId"]
	// Verify that the current user belongs in this list
	listUser := &models.ListUser{}
	if err := db.Where("list_id = ? AND user_id = ?", listID, userID).First(&listUser).Error; err != nil {
		return nil, err
	}

	item := &models.Item{
		ListID:   listUser.ListID,
		UserID:   userID,
		Name:     args["name"].(string),
		Quantity: args["quantity"].(int),
		Position: 1,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

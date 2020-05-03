package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// RetrieveItemsInList finds all items in a list by listID
func RetrieveItemsInList(db *gorm.DB, listID uuid.UUID) (interface{}, error) {
	items := []models.Item{}
	if err := db.Where("list_id = ?", listID).Order("completed ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

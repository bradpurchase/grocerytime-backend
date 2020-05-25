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
		Order("completed ASC").
		Order("created_at DESC").
		Find(&items).
		Error
	if err := query; err != nil {
		return nil, err
	}
	return items, nil
}

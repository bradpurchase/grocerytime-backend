package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// RetrieveItemsInList finds all items in a grocery trip by tripID
func RetrieveItemsInList(db *gorm.DB, tripID uuid.UUID) (interface{}, error) {
	items := []models.Item{}
	query := db.
		Where("grocery_trip_id = ?", tripID).
		Order("position ASC").
		Find(&items).
		Error
	if err := query; err != nil {
		return nil, err
	}
	return items, nil
}

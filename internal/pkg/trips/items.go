package trips

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// RetrieveItems finds all items in a grocery trip by tripID
func RetrieveItems(tripID uuid.UUID) (interface{}, error) {
	items := []models.Item{}
	query := db.Manager.
		Where("grocery_trip_id = ?", tripID).
		Order("position ASC").
		Find(&items).
		Error
	if err := query; err != nil {
		return nil, err
	}
	return items, nil
}

// RetrieveItemsInCategory finds all items in a grocery trip by category
func RetrieveItemsInCategory(tripID uuid.UUID, categoryID uuid.UUID) (interface{}, error) {
	items := []models.Item{}
	query := db.Manager.
		Where("grocery_trip_id = ?", tripID).
		Where("category_id = ?", categoryID).
		Order("position ASC").
		Find(&items).
		Error
	if err := query; err != nil {
		return nil, err
	}
	return items, nil
}

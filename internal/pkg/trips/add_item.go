package trips

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// AddItem adds an item to a trip and handles things like permission checks
func AddItem(db *gorm.DB, userID uuid.UUID, args map[string]interface{}) (interface{}, error) {
	tripID := args["tripId"]
	trip := &models.GroceryTrip{}
	if err := db.Where("id = ?", tripID).Find(&trip).Error; err != nil {
		return nil, err
	}

	// Verify that the current user belongs in this list
	listUser := &models.ListUser{}
	if err := db.Where("list_id = ? AND user_id = ?", trip.ListID, userID).First(&listUser).Error; err != nil {
		return nil, err
	}

	// TODO: support adding an item with category attributed
	// - Can we simply accept a categoryID?
	// - Do we accept a categoryName and then grab the ID from there?

	item := &models.Item{
		GroceryTripID: trip.ID,
		UserID:        userID,
		Name:          args["name"].(string),
		Quantity:      args["quantity"].(int),
		Position:      1,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

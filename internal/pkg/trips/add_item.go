package trips

import (
	"errors"

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

	itemCompleted := false
	item := models.Item{
		GroceryTripID: trip.ID,
		UserID:        userID,
		Name:          args["name"].(string),
		Quantity:      args["quantity"].(int),
		Position:      1,
		Completed:     &itemCompleted,
	}

	categoryName := args["categoryName"]
	category := &models.StoreCategory{}
	query := db.
		Select("store_categories.id").
		Joins("INNER JOIN categories ON categories.id = store_categories.category_id").
		Where("categories.name = ?", categoryName).
		First(&category).
		Error
	if err := query; err != nil {
		return nil, errors.New("could not find category")
	}
	item.CategoryID = &category.ID

	if err := db.Debug().Create(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

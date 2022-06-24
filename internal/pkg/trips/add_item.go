package trips

import (
	_ "embed"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// AddItem adds an item to a trip and handles things like permission checks
func AddItem(userID uuid.UUID, args map[string]interface{}) (addedItem *models.Item, err error) {
	tripID := args["tripId"].(uuid.UUID)

	itemCompleted := false
	itemName := args["name"].(string)
	quantity := 1
	if args["quantity"] != nil {
		quantity = args["quantity"].(int)
	}

	item := &models.Item{
		GroceryTripID: tripID,
		UserID:        userID,
		Name:          itemName,
		Quantity:      quantity,
		Position:      1,
		Completed:     &itemCompleted,
	}

	if args["stapleItemId"] != nil {
		stapleItemID := args["stapleItemId"].(uuid.UUID)
		item.StapleItemID = &stapleItemID
	}

	if err := db.Manager.Create(&item).Error; err != nil {
		return addedItem, err
	}
	return item, nil
}

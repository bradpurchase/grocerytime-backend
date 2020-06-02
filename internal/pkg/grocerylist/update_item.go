package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// UpdateItem updates an item by itemID
func UpdateItem(db *gorm.DB, args map[string]interface{}) (interface{}, error) {
	item := &models.Item{}
	if err := db.Where("id = ?", args["itemId"]).First(&item).Error; err != nil {
		return nil, err
	}

	if args["name"] != nil {
		item.Name = args["name"].(string)
	}
	if args["completed"] != nil {
		completed := args["completed"].(bool)
		item.Completed = completed
		item.Position = GetNewPosition(db, item.ListID, completed)
	}
	if args["quantity"] != nil {
		item.Quantity = args["quantity"].(int)
	}
	if args["position"] != nil {
		newPosition := args["position"].(int)
		currPosition := item.Position

		// Handle item position collisions
		if err := RepositionCollidingItem(db, item.ListID, currPosition, newPosition); err != nil {
			return nil, err
		}

		item.Position = newPosition
	}
	if err := db.Save(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// RepositionCollidingItem updates the position of a "colliding item" when performing
// an item update on position.
//
// When a user updates the position of an item in a list, the new position may collide
// with the position of an exising item. In this case, we need to update the colliding item's
// position either upward or downward in the list into a new position.
//
//TODO note: we may need to do this recursively; i.e. call this function within itself
// if the new position of the colliding item is also colliding
func RepositionCollidingItem(db *gorm.DB, listID uuid.UUID, currPosition int, newPosition int) error {
	collidingPosItem := &models.Item{}
	if err := db.Where("list_id = ? AND position = ?", listID, newPosition).Find(&collidingPosItem).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		return err
	}
	collidingItemNewPos := collidingPosItem.Position - 1
	// If the repositioned item's current position is greater than its new position,
	// move the colliding item's down the list by increasing its position by 1
	if currPosition > newPosition {
		collidingItemNewPos = collidingPosItem.Position + 1
	}
	if err := db.Model(&collidingPosItem).Where("id = ?", collidingPosItem.ID).Update("position", collidingItemNewPos).Error; err != nil {
		return err
	}
	return nil
}

// GetNewPosition gets the new position of an updated item
func GetNewPosition(db *gorm.DB, listID uuid.UUID, completed bool) int {
	// Reorder to the bottom of the list
	newPosition, err := DetermineListPosition("top", db, listID)
	if completed {
		newPosition, err = DetermineListPosition("bottom", db, listID)
	}
	if err != nil {
		return 0
	}
	return newPosition
}

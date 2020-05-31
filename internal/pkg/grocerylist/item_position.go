package grocerylist

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
)

// DetermineListPosition determines the position value to display
// an item at the top or bottom of a list depending on direction value
//
// Note: We position items by setting a default value of 1000, and decreasing
// the value by 2 each time an item is added. This is so we can efficiently
// order items; when we add a new item, we don't need to reorder all other items.
// We decrease by 2 so we can leave a gap for re-ordering, and when re-ordering
// we can do the same thing; we can just fit the re-ordered item within the "gap"
// between the two items around it without needing to touch them.
func DetermineListPosition(direction string, db *gorm.DB, listID interface{}) (int, error) {
	currPosItem := &models.Item{}
	orderDirection := "ASC"
	if direction == "bottom" {
		orderDirection = "DESC"
	}
	currPosItemQuery := db.
		Select("position").
		Where("list_id = ?", listID).
		Order("position " + orderDirection).
		Limit(1).
		Find(&currPosItem).
		Error
	if err := currPosItemQuery; err != nil && !gorm.IsRecordNotFoundError(err) {
		return 0, err
	}
	position := 1000
	// If there are no items yet set to 1000, otherwise add/subtract 2 depending on direction
	if currPosItem.Position > 0 {
		if direction == "bottom" {
			position = currPosItem.Position + 2
		} else {
			position = currPosItem.Position - 2
		}
	}
	return position, nil
}

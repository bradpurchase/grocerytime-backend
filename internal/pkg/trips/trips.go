package trips

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// RetrieveCurrentTripInList retrieves the currently active grocery trip in a
// list by listID if the userID has access to to the list
func RetrieveCurrentTripInList(db *gorm.DB, listID uuid.UUID, userID uuid.UUID) (interface{}, error) {
	if err := db.Where("list_id = ? AND user_id = ?", listID, userID).Find(&models.ListUser{}).Error; err != nil {
		return nil, errors.New("user is not a member of this list")
	}

	trip := models.GroceryTrip{}
	if err := db.Where("list_id = ? AND completed = ?", listID, false).Order("created_at DESC").Find(&trip).Error; err != nil {
		return nil, errors.New("could not find trip associated with this list")
	}
	return trip, nil
}

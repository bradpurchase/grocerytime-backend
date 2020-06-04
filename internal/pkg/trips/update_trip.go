package trips

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
)

// UpdateTrip updates a grocery trip with the given args by tripID
func UpdateTrip(db *gorm.DB, args map[string]interface{}) (interface{}, error) {
	trip := models.GroceryTrip{}
	if err := db.Where("id = ?", args["tripId"]).First(&trip).Error; err != nil {
		return nil, err
	}
	if args["name"] != nil {
		trip.Name = args["name"].(string)
	}
	if args["completed"] != nil {
		trip.Completed = args["completed"].(bool)
	}
	if err := db.Save(&trip).Error; err != nil {
		return nil, err
	}
	return trip, nil
}

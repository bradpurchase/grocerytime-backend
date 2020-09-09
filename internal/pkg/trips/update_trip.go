package trips

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
)

// UpdateTrip updates a grocery trip with the given args by tripID
func UpdateTrip(args map[string]interface{}) (interface{}, error) {
	trip := models.GroceryTrip{}
	if err := db.Manager.Where("id = ?", args["tripId"]).First(&trip).Error; err != nil {
		return nil, err
	}
	if args["name"] != nil {
		trip.Name = args["name"].(string)
	}
	if args["completed"] != nil {
		trip.Completed = args["completed"].(bool)
	}
	if args["copyRemainingItems"] != nil {
		trip.CopyRemainingItems = args["copyRemainingItems"].(bool)
	}
	if err := db.Manager.Save(&trip).Error; err != nil {
		return nil, err
	}
	return trip, nil
}

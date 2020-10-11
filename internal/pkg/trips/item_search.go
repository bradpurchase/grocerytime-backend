package trips

import (
	"errors"
	"fmt"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// SearchForItemByName finds an item within the stores the user belongs to
func SearchForItemByName(name string, userID uuid.UUID) (item *models.Item, err error) {
	var trips []models.GroceryTrip
	tripsQuery := db.Manager.
		Select("grocery_trips.id").
		Joins("INNER JOIN store_users ON grocery_trips.store_id = store_users.store_id").
		Where("store_users.user_id = ?", userID).
		Find(&trips).
		Error
	if err := tripsQuery; err != nil {
		return item, errors.New("could not find trip IDs for user")
	}
	var tripIDs []uuid.UUID
	for i := range trips {
		tripIDs = append(tripIDs, trips[i].ID)
	}

	foundItem := &models.Item{}
	nameArg := fmt.Sprintf("%%%s%%", name)
	searchQuery := db.Manager.
		Model(&models.Item{}).
		Where("name ILIKE ? AND grocery_trip_id IN ?", nameArg, tripIDs).
		First(&foundItem).
		Error
	if err := searchQuery; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return item, errors.New("no item matches the search term")
	}
	return foundItem, nil
}

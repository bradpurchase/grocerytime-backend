package trips

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"

	uuid "github.com/satori/go.uuid"
)

// RetrieveCurrentStoreTrip retrieves the currently active grocery trip in a
// store by storeID if the userID has access to to the store
func RetrieveCurrentStoreTrip(storeID uuid.UUID, user models.User) (groceryTrip models.GroceryTrip, err error) {
	query := db.Manager.
		Where("store_id = ?", storeID).
		Where("user_id = ? OR email = ?", user.ID, user.Email).
		Find(&models.StoreUser{}).
		Error
	if err := query; err != nil {
		return groceryTrip, errors.New("user is not a member of this store")
	}

	trip := models.GroceryTrip{}
	if err := db.Manager.Where("store_id = ? AND completed = ?", storeID, false).Order("created_at DESC").Find(&trip).Error; err != nil {
		return groceryTrip, errors.New("could not find trip associated with this store")
	}
	return trip, nil
}

// RetrieveTrips retrieves grocery trips within a store
func RetrieveTrips(storeID interface{}, userID uuid.UUID, completed bool) (tripsList []models.GroceryTrip, err error) {
	// Verify that this store belongs to the user
	var count int64
	existsQuery := db.Manager.
		Model(&models.StoreUser{}).
		Where("store_id = ? AND user_id = ? AND active = ?", storeID, userID, true).
		Count(&count).
		Error
	if err := existsQuery; err != nil {
		return tripsList, err
	}
	if count == 0 {
		return tripsList, errors.New("user is not active in this store")
	}

	var trips []models.GroceryTrip
	tripsQuery := db.Manager.
		Where("store_id = ? AND completed = ?", storeID, completed).
		Find(&trips).
		Error
	if err := tripsQuery; err != nil {
		return tripsList, err
	}
	return trips, nil
}

// RetrieveTrip retrieves a specific grocery trip by ID
func RetrieveTrip(tripID interface{}) (models.GroceryTrip, error) {
	trip := models.GroceryTrip{}

	if err := db.Manager.Where("id = ?", tripID).First(&trip).Error; err != nil {
		return trip, errors.New("trip not found")
	}
	return trip, nil
}

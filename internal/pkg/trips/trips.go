package trips

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// RetrieveCurrentStoreTrip retrieves the currently active grocery trip in a
// store by storeID if the userID has access to to the store
func RetrieveCurrentStoreTrip(db *gorm.DB, storeID uuid.UUID, user models.User) (interface{}, error) {
	query := db.
		Where("store_id = ?", storeID).
		Where("user_id = ? OR email = ?", user.ID, user.Email).
		Find(&models.StoreUser{}).
		Error
	if err := query; err != nil {
		return nil, errors.New("user is not a member of this store")
	}

	trip := models.GroceryTrip{}
	if err := db.Where("store_id = ? AND completed = ?", storeID, false).Order("created_at DESC").Find(&trip).Error; err != nil {
		return nil, errors.New("could not find trip associated with this store")
	}
	return trip, nil
}

// RetrieveTrip retrieves a specific grocery trip by ID
func RetrieveTrip(db *gorm.DB, tripID interface{}) (models.GroceryTrip, error) {
	trip := models.GroceryTrip{}

	if err := db.Where("id = ?", tripID).First(&trip).Error; err != nil {
		return trip, errors.New("trip not found")
	}
	return trip, nil
}

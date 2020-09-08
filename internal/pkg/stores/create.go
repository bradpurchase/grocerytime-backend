package stores

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// CreateStore creates a store for a user if it does not already exist by name
func CreateStore(db *gorm.DB, userID uuid.UUID, name string) (models.Store, error) {
	dupeStore, _ := RetrieveStoreForUserByName(db, name, userID)
	if dupeStore.Name != "" {
		return models.Store{}, errors.New("You already added a store with this name")
	}
	store := models.Store{UserID: userID, Name: name}
	if err := db.Debug().Create(&store).Error; err != nil {
		return models.Store{}, err
	}
	return store, nil
}

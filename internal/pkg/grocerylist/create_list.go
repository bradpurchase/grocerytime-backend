package grocerylist

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// CreateList creates a list for a user if it does not already exist by name
func CreateList(db *gorm.DB, userID uuid.UUID, name string) (models.List, error) {
	dupeList, _ := RetrieveListForUserByName(db, name, userID)
	if dupeList.Name != "" {
		return models.List{}, errors.New("You already have a list with this name")
	}
	list := models.List{UserID: userID, Name: name}
	if err := db.Create(&list).Error; err != nil {
		return models.List{}, err
	}
	return list, nil
}

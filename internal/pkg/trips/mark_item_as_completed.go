package trips

import (
	"errors"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// MarkItemAsCompleted item as completed by name for user (in any store)
func MarkItemAsCompleted(name string, userID uuid.UUID) (updatedItems []*models.Item, err error) {
	updateQuery := db.Manager.
		Model(&models.Item{}).
		Where("name = ? AND user_id = ?", name, userID).
		UpdateColumn("completed", true).
		Find(&updatedItems).
		Error
	if err := updateQuery; err != nil {
		return updatedItems, errors.New("could not update items")
	}
	return updatedItems, nil
}

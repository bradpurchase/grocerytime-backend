package notifications

import (
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	uuid "github.com/satori/go.uuid"
)

// StoreDeviceToken creates a Device record to store device tokens for iOS push notifications
func StoreDeviceToken(token string, userID uuid.UUID) (device *models.Device, err error) {
	userDevice := &models.Device{UserID: userID, Token: token}
	query := db.Manager.
		Where(userDevice).
		FirstOrCreate(&userDevice).
		Error
	if err := query; err != nil {
		return device, err
	}
	return userDevice, nil
}

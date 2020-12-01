package notifications

import (
	"errors"
	"fmt"
	"os"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db"
	uuid "github.com/satori/go.uuid"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

// Send sends a push notification
func Send(title string, body string, token string) (apnsID string, err error) {
	certFilename := "./certs/" + os.Getenv("APNS_CERT_FILENAME") + ".p12"
	cert, err := certificate.FromP12File(certFilename, os.Getenv("APNS_CERT_PASSWORD"))
	if err != nil {
		fmt.Printf("cert err: %v\n", err)
	}

	client := apns2.NewClient(cert).Development()

	notification := &apns2.Notification{}
	notification.DeviceToken = token
	payload := payload.NewPayload().AlertTitle(title).AlertBody(body)
	notification.Payload = payload

	res, err := client.Push(notification)
	if err != nil {
		fmt.Println("error:", err)
	}

	if res.Sent() {
		fmt.Println("[notifications/send] sent:", res.ApnsID)
		return res.ApnsID, nil
	}
	fmt.Printf("[notifications/send] not sent: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	return apnsID, errors.New(res.Reason)
}

// DeviceTokensForUser fetches all the device tokens stored for a user by ID
func DeviceTokensForUser(userID uuid.UUID) (tokens []string, err error) {
	query := db.Manager.
		Table("devices").
		Select("token").
		Where("user_id = ?", userID).
		Find(&tokens).
		Error
	if err := query; err != nil {
		return tokens, err
	}
	return tokens, nil
}

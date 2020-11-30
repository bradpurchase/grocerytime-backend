package notifications

import (
	"fmt"
	"os"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

// Send sends a push notification
func Send(message string, token string) {
	certFilename := "./certs/" + os.Getenv("APNS_CERT_FILENAME") + ".p12"
	cert, err := certificate.FromP12File(certFilename, os.Getenv("APNS_CERT_PASSWORD"))
	if err != nil {
		fmt.Printf("cert err: %v\n", err)
	}

	client := apns2.NewClient(cert).Development()

	notification := &apns2.Notification{}
	notification.DeviceToken = token
	payload := payload.NewPayload().Alert(message)
	notification.Payload = payload

	res, err := client.Push(notification)
	if err != nil {
		fmt.Println("error:", err)
	}

	if res.Sent() {
		fmt.Println("[notifications/send] notification sent:", res.ApnsID)
	} else {
		fmt.Printf("[notifications/send] notification not sent: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	}
}

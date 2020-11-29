package notifications

import (
	"fmt"
	"os"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

// Send sends a push notification
func Send() {
	cert, err := certificate.FromP12File("./certs/dev-cert.p12", os.Getenv("APNS_CERT_PASSWORD"))
	if err != nil {
		fmt.Printf("cert err: %v\n", err)
	}

	client := apns2.NewClient(cert).Development()

	notification := &apns2.Notification{}
	notification.DeviceToken = "e9c5c8cc94425b19a6a0126608fcb9e1ea5455101db2d79959e56b9305bc1f41"
	payload := payload.NewPayload().Alert("Test Notification")
	notification.Payload = payload

	res, err := client.Push(notification)
	if err != nil {
		fmt.Println("error:", err)
	}

	if res.Sent() {
		fmt.Println("Sent:", res.ApnsID)
	} else {
		fmt.Printf("Not Sent: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	}
}

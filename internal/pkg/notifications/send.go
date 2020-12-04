package notifications

import (
	"fmt"
	"os"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

// Send sends a push notification
func Send(title string, body string, token string, scheme string) {
	// Cert type and file path will change depending on environment (Debug, Beta, Release)
	certType := "prod"
	if scheme == "Debug" {
		certType = "test"
	}
	certFileName := fmt.Sprintf("%v-cert-%v", scheme, certType)
	certfile := fmt.Sprintf("certs/%v.p12", certFileName)
	cert, err := certificate.FromP12File(certfile, os.Getenv("APNS_CERT_PASSWORD"))
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
	} else {
		fmt.Printf("[notifications/send] not sent: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	}
}

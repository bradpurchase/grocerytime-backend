package notifications

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

// Send sends a push notification
func Send(title string, body string, token string, scheme string) {
	cert, err := ApnsCertificate(scheme)
	if err != nil {
		fmt.Printf("cert err: %v\n", err)
	}
	client := apns2.NewClient(cert).Production()
	if scheme == "Debug" {
		client = apns2.NewClient(cert).Development()
	}

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

// ApnsCertificate handles finding the proper certificate depending on app scheme and envrionment
func ApnsCertificate(scheme string) (cert tls.Certificate, err error) {
	certType := "prod"
	if scheme == "Debug" {
		certType = "test"
	}
	certFileName := fmt.Sprintf("%v-cert-%v", scheme, certType)
	certfile := fmt.Sprintf("%v/%v.p12", os.Getenv("APNS_CERT_FILEPATH"), certFileName)
	cert, err = certificate.FromP12File(certfile, os.Getenv("APNS_CERT_PASSWORD"))
	if err != nil {
		return cert, err
	}
	return cert, nil
}

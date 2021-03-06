package notifications

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

// Send sends a push notification
func Send(title string, body string, token string, entityName string, entityID string, scheme string) {
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
	notificationTopic := SetNotificationTopic(scheme)
	notification.Topic = notificationTopic
	notification.Payload = SetNotificationPayload(title, body, entityName, entityID)

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
	certPassword := ApnsCertificatePassword(scheme)
	cert, err = certificate.FromP12File(certfile, certPassword)
	if err != nil {
		return cert, err
	}
	return cert, nil
}

// ApnsCertificatePassword determines the credential for the
// apns cert password depending on environment
func ApnsCertificatePassword(scheme string) (password string) {
	cred := "APNS_CERT_PASSWORD"
	if scheme != "Debug" {
		cred = fmt.Sprintf("%v_%v", cred, strings.ToUpper(scheme))
	}
	return os.Getenv(cred)
}

// SetNotificationTopic determines the correct topic for the notification
// This should match the bundle ID depending on environment i.e. "bradpurchase.GroceryTime.beta"
func SetNotificationTopic(scheme string) (topic string) {
	topic = "bradpurchase.GroceryTime"
	if scheme != "Release" {
		topic = fmt.Sprintf("%v.%v", topic, strings.ToLower(scheme))
	}
	return topic
}

// SetNotificationPayload sets the APNS payload for the
func SetNotificationPayload(title string, body string, entityName string, entityID string) (p *payload.Payload) {
	return payload.
		NewPayload().
		AlertTitle(title).
		AlertBody(body).
		Custom("entity", entityName).
		Custom("id", entityID)
}

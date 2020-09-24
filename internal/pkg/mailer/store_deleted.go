package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendStoreDeletedEmail sends an email to a store user about a store being deleted
func SendStoreDeletedEmail(storeName string, userEmail string) (interface{}, error) {
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	subject := "Your list " + storeName + " has been deleted"
	to := mail.NewEmail("", userEmail)

	plainTextContent := "The list " + storeName + " has been deleted by the creator."
	htmlContent := "<p>Hello,</p>"
	htmlContent += "<p>This is to inform you that your list for the store " + storeName + " has been deleted. You can no longer access this store in the app.</p>"
	htmlContent += "<p>Thanks,<br />Brad from GroceryTime</p>"
	htmlContent += "<p>If you have any questions, concerns, or general feedback about GroceryTime, please email us at <a href=\"mailto:support@grocerytime.app\">support@grocerytime.app</a></p>"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	}
	return response, nil
}

package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendStoreRenamedEmail sends an email to a list user about a list being renamed
func SendStoreRenamedEmail(oldName string, newName string, email string) (interface{}, error) {
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	subject := "An update to one of your grocery lists in GroceryTime 📝"
	to := mail.NewEmail("", email)

	plainTextContent := "The list " + oldName + " was renamed to " + newName + "."
	htmlContent := "<p>Your list \"" + oldName + "\" has been renamed to \"" + newName + "\".</p>"
	htmlContent += "<p>----</p>"
	htmlContent += "<p>This message was sent to you because you are a member of the " + newName + " grocery list on GroceryTime.</p>"
	htmlContent += "<p>If you have any questions, concerns, or general feedback about GroceryTime, please email us at <a href=\"mailto:support@grocerytime.app\">support@grocerytime.app</a></p>"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	}
	return response, nil
}

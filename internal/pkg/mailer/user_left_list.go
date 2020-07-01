package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendUserLeftListEmail is sent to a list creator to inform them
// about a member of their list leaving
func SendUserLeftListEmail(listName string, listUserEmail string, recipientEmail string) (interface{}, error) {
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	subject := "Someone left your " + listName + " list 👋"
	to := mail.NewEmail("", recipientEmail)

	plainTextContent := "The member " + listUserEmail + " has left your list " + listName + "."

	htmlContent := "<p>Hello,</p>"
	htmlContent += "<p>A member of your list <strong>" + listName + "</strong>, "
	htmlContent += "<strong>" + listUserEmail + "</strong>, has left.</p>"
	htmlContent += "<p>This person no longer has access to this list. If you still "
	htmlContent += "want this person in your list, you are able to re-invite them by "
	htmlContent += "  tapping \"Share List\" in the list menu.</p>"
	htmlContent += "<p>Regards,<br />Brad from GroceryTime</p>"
	htmlContent += "<p>If you have any questions, concerns, or general feedback about GroceryTime, please email us at <a href=\"mailto:support@grocerytime.app\">support@grocerytime.app</a></p>"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	}
	return response, nil
}

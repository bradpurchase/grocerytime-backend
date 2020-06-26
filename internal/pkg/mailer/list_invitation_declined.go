package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendListInviteDeclinedEmail is sent to a list creator to inform them
// about a list invite to another user being declined
func SendListInviteDeclinedEmail(listName string, invitedUserEmail string, recipientEmail string) (interface{}, error) {
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	subject := "Invitation declined to join " + listName + " on GroceryTime 😞"
	to := mail.NewEmail("", recipientEmail)

	plainTextContent := "Your invitation sent to " + invitedUserEmail + " was declined"

	htmlContent := "<p>Hello,</p>"
	htmlContent += "<p>This email is to inform you that, sadly, your invitation "
	htmlContent += "sent to <strong>" + invitedUserEmail + "</strong> to join your "
	htmlContent += "list <strong>" + listName + "</strong> was declined.</p>"
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

package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendListInvitationEmail sends an email to a person being invited to join a list
func SendListInvitationEmail(listName string, email string) (interface{}, error) {
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	subject := "You've been invited to join " + listName + " on GroceryTime"
	to := mail.NewEmail("", email)

	plainTextContent := "You've been invited to join " + listName + " on GroceryTime. "
	plainTextContent += "Simply download the app and sign up with this email address to join. "
	plainTextContent += "Click here to download GroceryTime: https://grocerytime.app/download"

	htmlContent := "<p>Hello,</p>"
	htmlContent += "<p>You've been invited to join " + listName + ", a grocery list in the app GroceryTime.</p>"
	htmlContent += "<p>When you join a list that someone has created, you can work on it together and see updates in real time."
	htmlContent += "You can also split your list into as many separate trips as you'd like so you can track your grocery list from week to week (or any interval you like!)."
	htmlContent += "It makes grocery time so much easier!</p>"
	htmlContent += "<p>Simply download the app and sign up with this email address to join. Click here to download GroceryTime: https://grocerytime.app/download</p>"
	htmlContent += "<p>Thanks,<br />GroceryTime</p>"
	htmlContent += "<p>If you have any questions, concerns, or general feedback about GroceryTime, please email us at <a href=\"mailto:support@grocerytime.app\">support@grocerytime.app</a></p>"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	}
	return response, nil
}

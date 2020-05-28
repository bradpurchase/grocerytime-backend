package mailer

import (
	"os"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendListDeletedEmail sends an email to a list user about a list being deleted
func SendListDeletedEmail(list *models.List, user *models.User) (interface{}, error) {
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	subject := "Your list " + list.Name + " has been deleted"
	to := mail.NewEmail(user.FirstName+" "+user.LastName, user.Email)

	plainTextContent := "The list " + list.Name + " has been deleted by the creator."
	htmlContent := "<p>Hello,</p>"
	htmlContent += "<p>This is to inform you that the list " + list.Name + " has been deleted. You can no longer access this list or its items.</p>"
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

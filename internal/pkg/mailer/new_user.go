package mailer

import (
	"os"

	"github.com/bradpurchase/grocerytime-backend/internal/pkg/db/models"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendNewUserEmail sends a new
func SendNewUserEmail(user *models.User) (interface{}, error) {
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	subject := "Welcome to GroceryTime! 🛒"
	to := mail.NewEmail(user.FirstName+" "+user.LastName, user.Email)

	plainTextContent := "Hello, Welcome to GroceryTime, the collaborative grocery list app! This is just a quick email to thank you for signing up!"
	htmlContent := "<p>Hello,</p>"
	htmlContent += "<p>Welcome to GroceryTime, the collaborative grocery list app! This is just a quick email to thank you for signing up and give you a few pointers on how to get started.</p>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	}
	return response, nil
}

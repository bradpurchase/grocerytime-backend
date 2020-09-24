package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendNewUserEmail sends an email to a new user on signup
func SendNewUserEmail(email string) (interface{}, error) {
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	subject := "Welcome to GroceryTime! 🛒"
	to := mail.NewEmail("", email)

	plainTextContent := "Hello and welcome to GroceryTime, your grocery companion app! This is just a quick email to thank you for signing up."
	htmlContent := "<p>Hello,</p>"
	htmlContent += "<p>Welcome to GroceryTime, your grocery companion app! This is just a quick email to thank you for signing up and give you a few pointers on how to get started:</p>"
	htmlContent += "<ol>"
	htmlContent += "<li><strong>Create shopping lists:</strong> The first step after signing up is to tell us where you buy your groceries. This is so we can create a list for that store to get you started, but you can create as many as you'd like for each store you shop at!</li>"
	htmlContent += "<li><strong>Organize your shopping:</strong> Let's say you do a weekly grocery run. You can easily keep your list organized with separate trips to stay organized. Just tap the checkmark icon at the top of the list when you're done shopping!</li>"
	htmlContent += "<li><strong>Share with anyone:</strong> Most people who use GroceryTime use it to share grocery lists with their spouse or a friend. You can easily share any store list you've created and get work together - just use the share icon in the top right of the screen when viewing a store.</li>"
	htmlContent += "</ol>"
	htmlContent += "<p>Thanks again for signing up for a GroceryTime account. If you have any questions, concerns, or general feedback you can email us at <a href=\"mailto:support@grocerytime.app\">support@grocerytime.app</a>.</p>"
	htmlContent += "<p>Cheers,<br />Brad from GroceryTime</p>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	}
	return response, nil
}

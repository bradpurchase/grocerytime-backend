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

	plainTextContent := "Hello and welcome to GroceryTime, the collaborative grocery list app! This is just a quick email to thank you for signing up."
	htmlContent := "<p>Hello,</p>"
	htmlContent += "<p>Welcome to GroceryTime, the collaborative grocery list app! This is just a quick email to thank you for signing up and give you a few pointers on how to get started:</p>"
	htmlContent += "<ol>"
	htmlContent += "<li><strong>Download the app:</strong> You'll need the app on your phone, of course. Click here to download: <a href=\"https://groceryti.me/download\">https://groceryti.me/download</a></li>"
	htmlContent += "<li><strong>Create shopping lists:</strong> We've created a list called \"My Grocery List\" for you, but you can create as many lists as you want by tapping the \"+\" icon on the Lists screen.</li>"
	htmlContent += "<li><strong>Track grocery trips:</strong> Let's say you do a weekly grocery run. You can easily organize these weekly trips within your list, view your trip history, and much more.</li>"
	htmlContent += "<li><strong>List sharing:</strong> Most people who use GroceryTime use it to share grocery lists with their spouse or a friend. You can easily share any list with up to 5 people and get updates in real time - just tap the \":\" icon when viewing a list and then tap \"Share list\".</li>"
	htmlContent += "</ol>"
	htmlContent += "<p>Thanks again for signing up for a GroceryTime account. If you have any questions, concerns, or general feedback you can email us at <a href=\"mailto:support@grocerytime.app\">support@grocerytime.app</a>.</p>"
	htmlContent += "<p>Cheers,<br />Brad Purchase<br />Creator of GroceryTime</p>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	}
	return response, nil
}

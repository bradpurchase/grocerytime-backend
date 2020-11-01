package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendPasswordResetEmail sends an email to a user after their
// password was reset as part of the forgot password flow
func SendPasswordResetEmail(email string) (interface{}, error) {
	m := mail.NewV3Mail()
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	m.SetFrom(from)
	m.SetTemplateID("d-38da268588d24c84adad89bad41760b8")

	p := mail.NewPersonalization()
	toAddresses := []*mail.Email{
		mail.NewEmail("", email),
	}
	p.AddTos(toAddresses...)
	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

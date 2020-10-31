package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendNewUserEmail sends an email to a new user on signup
func SendNewUserEmail(email string) (interface{}, error) {
	m := mail.NewV3Mail()
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	m.SetFrom(from)
	m.SetTemplateID("d-323c957c7bdf42b095f96514e37db6e5")

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

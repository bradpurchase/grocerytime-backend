package mailer

import (
	"log"
	"os"

	uuid "github.com/satori/go.uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendForgotPasswordEmail sends an email to a user to reset their forgotten password
func SendForgotPasswordEmail(email string, token uuid.UUID) (interface{}, error) {
	m := mail.NewV3Mail()
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	m.SetFrom(from)
	m.SetTemplateID("d-ff98431efddb48619944cdba90406859")

	p := mail.NewPersonalization()
	toAddresses := []*mail.Email{
		mail.NewEmail("", email),
	}
	p.AddTos(toAddresses...)
	p.SetDynamicTemplateData("reset_token", token)
	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		return nil, err
	}
	log.Println(response)
	return response, nil
}

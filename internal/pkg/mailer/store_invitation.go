package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendStoreInvitationEmail sends an email to a person being invited to join a list
func SendStoreInvitationEmail(storeName string, email string, inviterName string) (interface{}, error) {
	m := mail.NewV3Mail()
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	m.SetFrom(from)
	m.SetTemplateID("d-39b0bd8d6b8747fcacbce147020364cd")

	p := mail.NewPersonalization()
	toAddresses := []*mail.Email{
		mail.NewEmail("", email),
	}
	p.AddTos(toAddresses...)

	p.SetDynamicTemplateData("first_name", inviterName)
	p.SetDynamicTemplateData("unique_name", storeName)
	p.SetDynamicTemplateData("email", email)

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

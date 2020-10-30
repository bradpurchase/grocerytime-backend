package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendStoreInviteDeclinedEmail is sent to a store creator to inform them
// about a store invite to another user being declined
func SendStoreInviteDeclinedEmail(storeName string, invitedUserEmail string, recipientEmail string) (interface{}, error) {
	m := mail.NewV3Mail()
	from := mail.NewEmail("GroceryTime", "noreply@groceryti.me")
	m.SetFrom(from)
	m.SetTemplateID("d-c7cf36bec4b94157adc9f4a90cd5e1c1")

	p := mail.NewPersonalization()
	toAddresses := []*mail.Email{
		mail.NewEmail("", recipientEmail),
	}
	p.AddTos(toAddresses...)

	p.SetDynamicTemplateData("email", invitedUserEmail)
	p.SetDynamicTemplateData("unique_name", storeName)

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

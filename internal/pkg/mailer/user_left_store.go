package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendUserLeftStoreEmail is sent to a list creator to inform them
// about a member of their list leaving
func SendUserLeftStoreEmail(storeName string, listUserName string, recipientEmail string) (interface{}, error) {
	m := mail.NewV3Mail()
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	m.SetFrom(from)
	m.SetTemplateID("d-63341b53948a4fee84da900aeae8f0f3")

	p := mail.NewPersonalization()
	toAddresses := []*mail.Email{
		mail.NewEmail("", recipientEmail),
	}
	p.AddTos(toAddresses...)

	p.SetDynamicTemplateData("user_name", listUserName)
	p.SetDynamicTemplateData("store_name", storeName)

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

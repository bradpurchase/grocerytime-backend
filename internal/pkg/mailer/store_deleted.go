package mailer

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendStoreDeletedEmail sends an email to a store member about a store being deleted
func SendStoreDeletedEmail(storeName string, emails []string) (interface{}, error) {
	m := mail.NewV3Mail()
	from := mail.NewEmail("GroceryTime", "noreply@grocerytime.app")
	m.SetFrom(from)
	m.SetTemplateID("d-523064949e1a4f739415896dacb80dc3")

	p := mail.NewPersonalization()
	// Retrieve emails for all StoreUser records
	toAddresses := []*mail.Email{}
	for i := range emails {
		toAddresses = append(toAddresses, mail.NewEmail("", emails[i]))
	}

	p.AddTos(toAddresses...)

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

/*
Implements an authboss.Mailer interface.
Use it with:
authboss.Config.Core.Mailer = authbossbrevomailer.NewMailer(authbossbrevomailer{})

The authboss.Email struct contains all the details for sending the email but
he email itself comes from the renderer. See renderer.go for details.

If we have no templateId defined in the config then this will throw an error
with the template name and parameters, which should make it easier to find the
correct name for the template map.
*/
package authbossbrevomailer

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/volatiletech/authboss/v3"
)

// Implements authboss Mailer interface.
type Mailer struct {
	Conf Config
}

func NewMailer(conf Config) Mailer {
	return Mailer{
		Conf: conf,
	}
}

// Lookup template ID from template map in the config.
func (m Mailer) GetTemplate(name string) int64 {
	val, ok := m.Conf.Template[name]
	if !ok {
		val = 0
	}
	return val
}

// Send rendered message using the Brevo API.
// https://developers.brevo.com/docs/send-a-transactional-email
func (m Mailer) Send(ctx context.Context, email authboss.Email) error {
	var data MailerData
	dec := gob.NewDecoder(bytes.NewReader([]byte(email.TextBody)))
	err := dec.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding email data", err)
		return err
	}

	// Lookup template ID from template map in the config.
	templateId := m.GetTemplate(data.Template)
	if templateId == 0 {
		err = fmt.Errorf("No template ID defined for enail %s", data.Template)
		fmt.Println(err)
		fmt.Println(data.Params)
		return err

	}

	// Populate recipients list.
	recipientsLen := len(email.To)
	if recipientsLen == 0 {
		return fmt.Errorf("No recipients for email %s", data.Template)
	}
	recipients := []BrevoRecipient{}
	recipientNamesLen := len(email.ToNames)
	for i := 0; i < recipientsLen; i++ {
		r := BrevoRecipient{Email: email.To[i]}
		if i < recipientNamesLen {
			r.Name = &email.ToNames[i]
		}
		recipients = append(recipients, r)
	}

	// We have the data from the renderer, send to API endpoint.
	payload := BrevoSendRequest{
		To:         recipients,
		TemplateId: templateId,
		Params:     data.Params,
	}

	if m.Conf.LogOnly {
		// Format a struct nicely for display.
		s, _ := json.MarshalIndent(payload, "", "\t")
		fmt.Printf("%s\n", s)
		return nil
	} else {

		b, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("Error marshalling email %s %s", data.Template, err)
		}
		req, err := http.NewRequest("POST", BREVO_API_URL, bytes.NewReader(b))
		if err != nil {
			return fmt.Errorf("Error creating email request %s %s", data.Template, err)
		}
		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")
		req.Header.Add("api-key", m.Conf.ApiKey)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("Error sending email request %s %s", data.Template, err)
		}
		if res.StatusCode != 201 {
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)
			err = fmt.Errorf("Status %d sending email %s %s", res.StatusCode, data.Template, body)
			fmt.Println(err)
			return err
		}
	}

	return nil

}

/*
Implements an authboss.Renderer interface that prepares an authboss.Email for
sending with authbossbrevomailer.Mailer.

Use with: authboss.Config.Core.MailRenderer = authbossbrevomailer.Renderer{}

Authboss assumes that the email will be send via SMTP to created text baased
email from templates. As the Brevo HTTP API expects a JSON object we use the
mechanism provided for text email content to pass the encoded details from the
renderer here, to the Mailer.
*/

package authbossbrevomailer

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/volatiletech/authboss/v3"
)

const SUFFIX string = "_txt"

type MailerData struct {
	Template string
	Params   authboss.HTMLData
}

type Renderer struct{}

// Load template.
// As we are using Brevo for rendering, we do not need to load any templates.
func (r Renderer) Load(names ...string) error {
	return nil
}

// Render the template.
// As we are using the Brevo API, we do not need to render any templates.
// instead merge the template name and data and encode for Send.
func (r Renderer) Render(ctx context.Context, page string, data authboss.HTMLData) ([]byte, string, error) {
	b := []byte{}
	// Render is called twice, once for HTML and once for text. We only want to
	// do this once, so encode the data in the text channel.
	if strings.HasSuffix(page, SUFFIX) {
		delete(data, "modules")

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(MailerData{
			Template: page[:len(page)-len(SUFFIX)],
			Params:   data,
		})
		if err == nil {
			b = buf.Bytes()
		} else {
			fmt.Println("Error encoding data", err)
		}
	}

	// No need to render a template, return nothing.
	return b, "", nil
}

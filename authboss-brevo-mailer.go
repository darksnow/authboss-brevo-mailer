/*
Package authbossbrevomailer is an authboss Mailer for sending transactional
emails via Brevo (https://brevo.com) using v3 of the HTTP API.
*/
package authbossbrevomailer

// The API endpoint used to send emails.
const BREVO_API_URL string = "https://api.brevo.com/v3/smtp/email"

// Configuration object. Typically populated from envvars.
type Config struct {
	// API Key from Brevo.
	ApiKey string
	// Set to false to actually send email.
	LogOnly bool
	// Map of template names and associated Brevo template ID.
	Template map[string]int64
}

// Structure of a request to the Brevo API.
// https://developers.brevo.com/reference/sendtransacemail
type BrevoRecipient struct {
	Email string  `json:"email"`
	Name  *string `json:"name,omitempty"`
}

// Sender is not required when using template IDs.
type BrevoSendRequest struct {
	To         []BrevoRecipient       `json:"to,omitempty"`         // Email address and optional name of recipient.
	TemplateId int64                  `json:"templateId,omitempty"` // Brevo TemplateID for the email template we want to send.
	Params     map[string]interface{} `json:"params,omitempty"`     // List of paramaters defined in the template.
}

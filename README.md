# Authboss Brevo Mailer

An [Authboss](https://github.com/volatiletech/authboss) Mailer and Renderer to send transactional emails with [Brevo](https://www.brevo.com/).

## Usage

You will need a Brevo account. Then you can get an API key from https://app.brevo.com/settings/keys/smtp

Once you create a template in Brevo, you can assign the ID of that template to the configuration which maps the name of the Authboss template to the Brevo ID of the template you want to send.

As Brevo templates contain the sender address along with all details needed to build the email, all you will need is this ID.

Sample code for integrating this mailer can be seen below:

    conf := authbossbrevomailer.Config{
        // API Key from Brevo.
        ApiKey: "API_KEY_FROM_BREVO",
        // Set to false to actually send email.
        LogOnly: false,
        // Map of template names and associated Brevo template ID.
        Template: map[string]int64 {
            "recover": 5,
        },
    }

    ab = authboss.New()
    authboss.Config.Core.Mailer = authbossbrevomailer.NewMailer(conf)
    authboss.Config.Core.MailRenderer = authbossbrevomailer.Renderer{}

**Both the Mailer and MailRenderer MUST be set for this to work** as the Mailer requires the template parameters to be encoded by the renderer to work. See Methodology below for details.

## Templates

If no template ID can be found in the configuration then the sender  will throw an error after logging the template name and parameters. This should make it easier to find the correct name for the template map and the parameters needed in your Brevo template.

### Authboss templates and parameters

These are the templates and parameters taken from the default Authboss renderer email templates. https://github.com/volatiletech/authboss-renderer/tree/master/email-templates

| Template name | Paramaters | Purpose |
| --- | --- | --- |
| `confirm` | `url` | URL containing key for email confirmation. |
| `recover` | `recover_url` | URL containing key for password recovery. |
| `twofactor_verify_email` | `url` | URL containing key for 2FA verification. |

Note that the parameters listed above should be prefixed with `params` before being added to the Brevo template. i.e. `{{ params.url }}`. See the [Brevo developer guide](https://developers.brevo.com/docs/send-a-transactional-email) for more details.

You only need to implement the templates you are actually using so if you have not enabled 2FA, for example, then you do not need to create the template in Brevo and set the ID in the configuration.

## Methodology

The `authboss.Email` struct contains all the details for sending the email but the email parameters come from the renderer.

In `renderer.go` the code encodes the parameters and template name into the body of a plain text email. Later, in the mailer, those details are recovered so they can be passed to the API.

Authboss is setup to render HTML and plain text parts of a MIME email in the renderer. These are then assumed to be the body of an email which can be passed to an SMTP server in the mailer. Since we want to send using a JSON API endpoint we need the parameters in a more structured format. As an updated GO `context` was not available and the parameters are not directly accessible from the mailer, this was the solution I came up with.

Brevo also has a full featured [GO SDK](https://github.com/getbrevo/brevo-go) for accessing the API. As all we need to do is send an email, which can be done with an HTTP endpoint, there was no need to add a dependency on the full library.

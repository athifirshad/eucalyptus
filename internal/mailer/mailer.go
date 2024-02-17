package mailer

import (
	"bytes"
	"embed"
	"text/template"

	"github.com/wneessen/go-mail"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	client *mail.Client
	sender string
}

func New(host string, port int, username, password, sender string) (*Mailer, error) {
	// Configure the mail client with the provided SMTP details
	client, err := mail.NewClient(host, mail.WithPort(port), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(username), mail.WithPassword(password))
	if err != nil {
		return nil, err
	}

	// If the server supports TLS, enable it
	if port ==  465 {
		client.SetSSLPort(true, false)
	} else if port ==  587 {
		client.SetTLSPortPolicy(mail.TLSMandatory)
	}

	return &Mailer{
		client: client,
		sender: sender,
	}, nil
}


func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}
	msg := mail.NewMsg()
	msg.To(recipient)
	msg.From(m.sender)
	msg.Subject( subject.String())
	msg.SetBodyString("text/plain", plainBody.String())
	msg.AddAlternativeString("text/html", htmlBody.String())
	err = m.client.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}

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

func New(host string, port int, username, password, sender string) Mailer {
	client, _ := mail.NewClient(host, mail.WithPort(port), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(username), mail.WithPassword(password), mail.WithSSLPort(true))
	defer client.Close()
	return Mailer{
		client: client,
		sender: sender,
	}
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
	msg.SetGenHeader("To", recipient)
	msg.SetGenHeader("From", m.sender)
	msg.SetGenHeader("Subject", subject.String())
	msg.SetBodyString("text/plain", plainBody.String())
	msg.AddAlternativeString("text/html", htmlBody.String())
	err = m.client.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}

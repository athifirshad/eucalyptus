package mailer

import (
	"bytes"
	"embed"
	"text/template"

	"github.com/athifirshad/eucalyptus/internal/tasks"
	"github.com/hibiken/asynq"
	"github.com/wneessen/go-mail"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	client      *mail.Client
	sender      string
	asynqClient *asynq.Client
}

func NewMailer(host string, port int, username, password, sender string, asynqClient *asynq.Client) (*Mailer, error) {
	client, err := mail.NewClient(host, mail.WithPort(port), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(username), mail.WithPassword(password))
	if err != nil {
		return nil, err
	}

	client.SetTLSPortPolicy(mail.TLSMandatory)
	return &Mailer{
		client:      client,
		sender:      sender,
		asynqClient: asynqClient,
	}, nil
}

func (m Mailer) Send(recipient, templateFile string, data map[string]any) error {
	task, _ := tasks.NewEmailTask(recipient, templateFile, data)

	_, err := m.asynqClient.Enqueue(task, asynq.MaxRetry(5))
	if err != nil {
		return err
	}
	return err
}

func (m Mailer) SendEmail(recipient, templateFile string, data map[string]any) error {
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
	msg.Subject(subject.String())
	msg.SetBodyString("text/plain", plainBody.String())
	msg.AddAlternativeString("text/html", htmlBody.String())
	err = m.client.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}

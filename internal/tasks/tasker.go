package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

type EmailTask struct {
	Recipient string
	Template  string
	Data      map[string]interface{}
}

func (t *EmailTask) Type() string {
	return "email:send"
}
func (t *EmailTask) Payload() ([]byte, error) {
	return json.Marshal(t)
}
func NewEmailTask(recipient, template string, data map[string]any) (*asynq.Task,error) {
	payload, err := json.Marshal(EmailTask{
		Recipient: recipient,
		Template:  template,
		Data:      data,
	})
	if err != nil {
		return nil,err

	}
	return asynq.NewTask("email:send", payload),nil
}

package mail

import (
	"bytes"
	"errors"
	"html/template"

	"chefbook-server/internal/config"
	"chefbook-server/pkg/logger"
)

type Sender interface {
	Send(input SendEmailInput) error
}

func NewSender(config config.SMTPConfig, useSMTP bool) (Sender, error) {
	if useSMTP {
		return NewSMTPSender(config.From, config.Password, config.Host, config.Port)
	}
	return NewFakeSMTP(), nil
}

type SendEmailInput struct {
	To      string
	Subject string
	Body    string
}

func (e *SendEmailInput) GenerateBodyFromHTML(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		logger.Errorf("failed to parse file %s:%s", templateFileName, err.Error())

		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	e.Body = buf.String()

	return nil
}

func (e *SendEmailInput) Validate() error {
	if e.To == "" {
		return errors.New("empty to")
	}

	if e.Subject == "" || e.Body == "" {
		return errors.New("empty subject/body")
	}

	if !IsEmailValid(e.To) {
		return errors.New("invalid to email")
	}

	return nil
}

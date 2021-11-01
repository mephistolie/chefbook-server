package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/config"
	"github.com/mephistolie/chefbook-server/pkg/cache"
	emailProvider "github.com/mephistolie/chefbook-server/pkg/mail"
)

const (
	verificationLinkTmpl = "https://%s/v1/users/activate/%s"
)

type EmailService struct {
	sender emailProvider.Sender
	config config.MailConfig

	cache cache.Cache
}

type verificationEmailInput struct {
	VerificationLink string
}

func NewMailService(sender emailProvider.Sender, config config.MailConfig, cache cache.Cache) *EmailService {
	return &EmailService{
		sender: sender,
		config: config,
		cache:  cache,
	}
}

func (s *EmailService) SendVerificationEmail(input VerificationEmailInput) error {
	subject := fmt.Sprint(s.config.Subjects.Verification)

	templateInput := verificationEmailInput{s.createVerificationLink(input.Domain, input.VerificationCode)}
	sendInput := emailProvider.SendEmailInput{Subject: subject, To: input.Email}

	if err := sendInput.GenerateBodyFromHTML(s.config.Templates.Verification, templateInput); err != nil {
		return err
	}

	return s.sender.Send(sendInput)
}

func (s *EmailService) createVerificationLink(domain string, code uuid.UUID) string {
	return fmt.Sprintf(verificationLinkTmpl, domain, code)
}

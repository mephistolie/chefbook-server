package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/config"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/pkg/cache"
	emailProvider "github.com/mephistolie/chefbook-server/pkg/mail"
)

const (
	verificationLinkTmpl = "https://%s/v1/auth/activate/%s"
)

type MailService struct {
	sender emailProvider.Sender
	config config.MailConfig

	cache cache.Cache
}

type verificationEmailInput struct {
	VerificationLink string
}

func NewMailService(sender emailProvider.Sender, config config.MailConfig, cache cache.Cache) *MailService {
	return &MailService{
		sender: sender,
		config: config,
		cache:  cache,
	}
}

func (s *MailService) SendVerificationEmail(input entity.VerificationEmailInput) error {
	subject := fmt.Sprint(s.config.Subjects.Verification)

	templateInput := verificationEmailInput{s.createVerificationLink(input.Domain, input.VerificationCode)}
	sendInput := emailProvider.SendEmailInput{Subject: subject, To: input.Email}

	if err := sendInput.GenerateBodyFromHTML(s.config.Templates.Verification, templateInput); err != nil {
		return err
	}

	if err := s.sender.Send(sendInput); err != nil {
		return failure.UnableSendEmail
	}

	return nil
}

func (s *MailService) createVerificationLink(domain string, code uuid.UUID) string {
	return fmt.Sprintf(verificationLinkTmpl, domain, code)
}

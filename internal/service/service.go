package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/config"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"github.com/mephistolie/chefbook-server/pkg/auth"
	"github.com/mephistolie/chefbook-server/pkg/cache"
	"github.com/mephistolie/chefbook-server/pkg/hash"
	"github.com/mephistolie/chefbook-server/pkg/mail"
	"time"
)

type Users interface {
	SignUp(authInput models.AuthData) (int, error)
	ActivateUser(activationLink uuid.UUID) error
	SignIn(authInput models.AuthData, ip string) (models.Tokens, error)
	RefreshSession(refreshToken, ip string) (models.Tokens, error)
}

type VerificationEmailInput struct {
	Email            string
	Name             string
	VerificationCode uuid.UUID
	Domain           string
}

type Mails interface {
	SendVerificationEmail(input VerificationEmailInput) error
}

type Recipes interface {
	AddRecipe(recipe models.Recipe) (int, error)
}

type Service struct {
	Users
	Mails
	Recipes
}

type Dependencies struct {
	Repos            *repository.Repository
	Cache            cache.Cache
	HashManager      hash.HashManager
	TokenManager     auth.TokenManager
	MailSender       mail.Sender
	MailConfig       config.MailConfig
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	FondyCallbackURL string
	CacheTTL         int64
	Environment      string
	Domain           string
}

func NewServices(dependencies Dependencies) *Service {

	mailService := NewMailService(dependencies.MailSender, dependencies.MailConfig, dependencies.Cache)

	return &Service{
		Users: NewUsersService(dependencies.Repos, dependencies.HashManager, dependencies.TokenManager,
			dependencies.AccessTokenTTL, dependencies.RefreshTokenTTL, mailService),
		Mails: NewMailService(dependencies.MailSender, dependencies.MailConfig, dependencies.Cache),
		Recipes: NewRecipesService(dependencies.Repos),
	}
}

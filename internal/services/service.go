package services

import (
	"github.com/mephistolie/chefbook-server/internal/models"
	services "github.com/mephistolie/chefbook-server/internal/repositories"
	"github.com/spf13/viper"
	"os"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(email, password string) (string, error)
}

type Mail interface {
	SendEmailVerificationCode(code int, recipient string) error
}

type Recipes interface {

}

type Service struct {
	Authorization
	Mail
	Recipes
}

func NewService(repos *services.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Mail: NewMailService(viper.GetString("smtp.host"), viper.GetString("smtp.port"),
			os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD")),
	}
}
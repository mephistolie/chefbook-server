package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"github.com/mephistolie/chefbook-server/pkg/hash"
	"time"
)

const (
	signingKey = ""
	tokenTTL = 12 * time.Hour
)

type UsersService struct {
	repo        repository.Users
	mailService Mails
	hashManager hash.HashManager
}

func NewAuthService(repo repository.Users, hashManager hash.HashManager, mailService Mails) *UsersService {
	return &UsersService{
		repo: repo,
		hashManager: hashManager,
		mailService: mailService,
	}
}

func (s *UsersService) CreateUser(user models.User) (int, error) {
	if candidate, _ := s.repo.GetUserByEmail(user.Email); candidate.Id != 0 && candidate.IsActivated == false {
		err := s.mailService.SendVerificationEmail(VerificationEmailInput{
			Email: user.Email,
			Domain: "localhost",
			VerificationCode: candidate.ActivationLink,
		})
		return candidate.Id, err
	} else if candidate.Id != 0 {
		return 0, errors.New("user already registered")
	}

	hashedPassword, err := s.hashManager.Hash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword
	activationLink := uuid.New()
	userId, err := s.repo.CreateUser(user, activationLink);
	if err != nil {
		return 0, err
	}
	err = s.mailService.SendVerificationEmail(VerificationEmailInput{
		Email: user.Email,
		Domain: "localhost",
		VerificationCode: activationLink,
	})
	if err != nil {
		return 0, err
	}
	return userId, err
}

func (s *UsersService) ActivateUser(activationLink uuid.UUID) error {
	return s.repo.ActivateUser(activationLink)
}
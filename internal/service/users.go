package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"github.com/mephistolie/chefbook-server/pkg/auth"
	"github.com/mephistolie/chefbook-server/pkg/hash"
	"strconv"
	"time"
)

type UsersService struct {
	repo repository.Users

	hashManager hash.HashManager

	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration

	mailService Mails
}

func NewUsersService(repo repository.Users, hashManager hash.HashManager, tokenManager auth.TokenManager,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration, mailService Mails) *UsersService {
	return &UsersService{
		repo:            repo,
		hashManager:     hashManager,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		mailService:     mailService,
	}
}

func (s *UsersService) SignUp(authData models.AuthData) (int, error) {
	hashedPassword, err := s.hashManager.Hash(authData.Password)
	if err != nil {
		return 0, err
	}

	if candidate, _ := s.repo.GetUserByEmail(authData.Email); candidate.Id != 0 && candidate.IsActivated == false {
		if candidate.Password != hashedPassword {
			s.repo.ChangePassword(authData)
		}
		err := s.mailService.SendVerificationEmail(VerificationEmailInput{
			Email:            authData.Email,
			Domain:           "localhost:8000",
			VerificationCode: candidate.ActivationLink,
		})
		return candidate.Id, err
	} else if candidate.Id != 0 {
		return 0, models.ErrUserAlreadyExists
	}

	authData.Password = hashedPassword
	activationLink := uuid.New()
	userId, err := s.repo.CreateUser(authData, activationLink)
	if err != nil {
		return 0, err
	}
	err = s.mailService.SendVerificationEmail(VerificationEmailInput{
		Email:            authData.Email,
		Domain:           "localhost:8000",
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

func (s *UsersService) SignIn(authData models.AuthData, ip string) (models.Tokens, error) {
	user, err := s.repo.GetUserByEmail(authData.Email)
	if err != nil {
		return models.Tokens{}, models.ErrUserNotFound
	}
	if user.IsActivated == false {
		return models.Tokens{}, models.ErrProfileNotActivated
	}
	if user.IsBlocked == true {
		return models.Tokens{}, models.ErrProfileIsBlocked
	}
	if err = s.hashManager.ValidateByHash(authData.Password, user.Password); err != nil {
		return models.Tokens{}, models.ErrAuthentication
	}
	return s.CreateSession(user.Id, ip)
}

func (s *UsersService) CreateSession(userId int, ip string) (models.Tokens, error) {
	var (
		res models.Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(strconv.Itoa(userId), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := models.Session{
		UserId: userId,
		RefreshToken: res.RefreshToken,
		Ip: ip,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	return res, s.repo.CreateSession(session)
}

func (s *UsersService) RefreshSession(currentRefreshToken, ip string) (models.Tokens, error) {
	user, err := s.repo.GetByRefreshToken(currentRefreshToken)
	if err != nil {
		return models.Tokens{}, err
	}
	if user.IsBlocked == true {
		if err := s.repo.DeleteSession(currentRefreshToken); err != nil {
			return models.Tokens{}, err
		}
		return models.Tokens{}, models.ErrProfileIsBlocked
	}

	var res models.Tokens

	res.AccessToken, err = s.tokenManager.NewJWT(strconv.Itoa(user.Id), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := models.Session{
		UserId: user.Id,
		RefreshToken: res.RefreshToken,
		Ip: ip,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	return res, s.repo.UpdateSession(session, currentRefreshToken)
}
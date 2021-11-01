package service

import (
	"github.com/google/uuid"
	models2 "github.com/mephistolie/chefbook-server/internal/models"
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
	domain      string
}

func NewUsersService(repo repository.Users, hashManager hash.HashManager, tokenManager auth.TokenManager,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration, mailService Mails, domain string) *UsersService {
	return &UsersService{
		repo:            repo,
		hashManager:     hashManager,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		mailService:     mailService,
		domain:          domain,
	}
}

func (s *UsersService) SignUp(authData models2.AuthData) (int, error) {
	hashedPassword, err := s.hashManager.Hash(authData.Password)
	if err != nil {
		return -1, err
	}

	if candidate, _ := s.repo.GetUserByEmail(authData.Email); candidate.Id > 0 && candidate.IsActivated == false {
		if err = s.hashManager.ValidateByHash(authData.Password, candidate.Password); err != nil {
			authData.Password = hashedPassword
			err := s.repo.ChangePassword(authData)
			if err != nil {
				return -1, err
			}
		}
		if candidate.Password != hashedPassword {

		}
		err := s.mailService.SendVerificationEmail(VerificationEmailInput{
			Email:            authData.Email,
			Domain:           s.domain,
			VerificationCode: candidate.ActivationLink,
		})
		return candidate.Id, err
	} else if candidate.Id > 0 {
		return 0, models2.ErrUserAlreadyExists
	}

	authData.Password = hashedPassword
	activationLink := uuid.New()
	userId, err := s.repo.CreateUser(authData, activationLink)
	if err != nil {
		return 0, err
	}
	err = s.mailService.SendVerificationEmail(VerificationEmailInput{
		Email:            authData.Email,
		Domain:           s.domain,
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

func (s *UsersService) SignIn(authData models2.AuthData, ip string) (models2.Tokens, error) {
	user, err := s.repo.GetUserByEmail(authData.Email)
	if err != nil {
		return models2.Tokens{}, models2.ErrUserNotFound
	}
	if user.IsActivated == false {
		return models2.Tokens{}, models2.ErrProfileNotActivated
	}
	if user.IsBlocked == true {
		return models2.Tokens{}, models2.ErrProfileIsBlocked
	}
	if err = s.hashManager.ValidateByHash(authData.Password, user.Password); err != nil {
		return models2.Tokens{}, models2.ErrAuthentication
	}
	return s.CreateSession(user.Id, ip)
}

func (s *UsersService) CreateSession(userId int, ip string) (models2.Tokens, error) {
	var (
		res models2.Tokens
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

	session := models2.Session{
		UserId:       userId,
		RefreshToken: res.RefreshToken,
		Ip:           ip,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	return res, s.repo.CreateSession(session)
}

func (s *UsersService) RefreshSession(currentRefreshToken, ip string) (models2.Tokens, error) {
	user, err := s.repo.GetByRefreshToken(currentRefreshToken)
	if err != nil {
		return models2.Tokens{}, err
	}
	if user.IsBlocked == true {
		if err := s.repo.DeleteSession(currentRefreshToken); err != nil {
			return models2.Tokens{}, err
		}
		return models2.Tokens{}, models2.ErrProfileIsBlocked
	}

	var res models2.Tokens

	res.AccessToken, err = s.tokenManager.NewJWT(strconv.Itoa(user.Id), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := models2.Session{
		UserId:       user.Id,
		RefreshToken: res.RefreshToken,
		Ip:           ip,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	return res, s.repo.UpdateSession(session, currentRefreshToken)
}

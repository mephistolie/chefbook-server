package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/model"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"github.com/mephistolie/chefbook-server/pkg/auth"
	"github.com/mephistolie/chefbook-server/pkg/hash"
	"strconv"
	"time"
)

type AuthService struct {
	repo repository.Auth

	hashManager     hash.HashManager
	firebaseService FirebaseService

	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration

	mailService Mails
	domain      string
}

func NewAuthService(repo repository.Auth, firebaseService FirebaseService, hashManager hash.HashManager, tokenManager auth.TokenManager,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration, mailService Mails, domain string) *AuthService {
	return &AuthService{
		repo:            repo,
		firebaseService: firebaseService,
		hashManager:     hashManager,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		mailService:     mailService,
		domain:          domain,
	}
}

func (s *AuthService) SignUp(authData model.AuthData) (int, error) {
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
		err := s.mailService.SendVerificationEmail(VerificationEmailInput{
			Email:            authData.Email,
			Domain:           s.domain,
			VerificationCode: candidate.ActivationLink,
		})
		return candidate.Id, err
	} else if candidate.Id > 0 {
		return 0, model.ErrUserAlreadyExists
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

func (s *AuthService) ActivateUser(activationLink uuid.UUID) error {
	return s.repo.ActivateUser(activationLink)
}

func (s *AuthService) SignIn(authData model.AuthData, ip string) (model.Tokens, error) {
	user, err := s.repo.GetUserByEmail(authData.Email)
	password := authData.Password
	if err != nil {
		firebaseUser, err := s.firebaseService.SignIn(authData)
		if err != nil {
			return model.Tokens{}, model.ErrUserNotFound
		}
		hashedPassword, err := s.hashManager.Hash(authData.Password)
		if err != nil {
			return model.Tokens{}, err
		}
		authData.Password = hashedPassword
		err = s.firebaseService.migrateFromFirebase(authData, firebaseUser)
		if err != nil {
			return model.Tokens{}, err
		}
		user, err = s.repo.GetUserByEmail(authData.Email)
		if err != nil {
			return model.Tokens{}, err
		}
	}
	if user.IsActivated == false {
		return model.Tokens{}, model.ErrProfileNotActivated
	}
	if user.IsBlocked == true {
		return model.Tokens{}, model.ErrProfileIsBlocked
	}
	if err = s.hashManager.ValidateByHash(password, user.Password); err != nil {
		return model.Tokens{}, model.ErrAuthentication
	}
	return s.CreateSession(user.Id, ip)
}

func (s *AuthService) CreateSession(userId int, ip string) (model.Tokens, error) {
	var (
		res model.Tokens
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

	session := model.Session{
		UserId:       userId,
		RefreshToken: res.RefreshToken,
		Ip:           ip,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	return res, s.repo.CreateSession(session)
}

func (s *AuthService) SignOut(refreshToken string) error {
	if err := s.repo.DeleteSession(refreshToken); err != nil {
		return model.ErrSessionNotFound
	}
	return nil
}

func (s *AuthService) RefreshSession(currentRefreshToken, ip string) (model.Tokens, error) {
	user, err := s.repo.GetByRefreshToken(currentRefreshToken)
	if err != nil {
		return model.Tokens{}, model.ErrSessionNotFound
	}
	if user.IsBlocked == true {
		if err := s.repo.DeleteSession(currentRefreshToken); err != nil {
			return model.Tokens{}, err
		}
		return model.Tokens{}, model.ErrProfileIsBlocked
	}

	var res model.Tokens

	res.AccessToken, err = s.tokenManager.NewJWT(strconv.Itoa(user.Id), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := model.Session{
		UserId:       user.Id,
		RefreshToken: res.RefreshToken,
		Ip:           ip,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	return res, s.repo.UpdateSession(session, currentRefreshToken)
}

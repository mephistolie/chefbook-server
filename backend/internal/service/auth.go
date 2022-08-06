package service

import (
	"chefbook-server/internal/entity"
	"chefbook-server/internal/entity/failure"
	"chefbook-server/internal/service/interface/repository"
	"chefbook-server/pkg/auth"
	"chefbook-server/pkg/hash"
	"github.com/google/uuid"
	"strconv"
	"time"
)

const maxSessionsCount = 5

type AuthService struct {
	repo repository.Auth

	hashManager     hash.HashManager
	firebaseService *FirebaseService

	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration

	mailService MailService
	domain      string
}

func NewAuthService(repo repository.Auth, firebaseService *FirebaseService, hashManager hash.HashManager, tokenManager auth.TokenManager,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration, mailService MailService, domain string) *AuthService {
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

func (s *AuthService) SignUp(credentials entity.Credentials) (int, error) {
	hashedPassword, err := s.hashManager.Hash(credentials.Password)
	if err != nil {
		return 0, failure.Unknown
	}

	if candidate, err := s.repo.GetUserByEmail(credentials.Email); err == nil {
		if candidate.IsActivated {
			return 0, failure.UserAlreadyExists
		}

		if err = s.hashManager.ValidateByHash(credentials.Password, candidate.Password); err != nil {
			credentials.Password = hashedPassword
			err := s.repo.ChangePassword(candidate.Id, credentials.Password)
			if err != nil {
				return 0, failure.Unknown
			}
		}

		activationLink, err := s.repo.GetUserActivationLink(candidate.Id)
		if err != nil {
			return 0, failure.Unknown
		}
		return candidate.Id, s.sendActivationLink(credentials.Email, activationLink)
	}

	credentials.Password = hashedPassword
	activationLink := uuid.New()
	userId, err := s.repo.CreateUser(credentials, activationLink)
	if err != nil {
		return 0, err
	}

	return userId, s.sendActivationLink(credentials.Email, activationLink)
}

func (s *AuthService) ActivateProfile(activationLink uuid.UUID) error {
	return s.repo.ActivateProfile(activationLink)
}

func (s *AuthService) SignIn(credentials entity.Credentials, ip string) (entity.Tokens, error) {
	user, err := s.repo.GetUserByEmail(credentials.Email)
	password := credentials.Password

	if err != nil && s.firebaseService != nil {
		if migratedUser, err := s.migrateFromFirebase(credentials); err != nil {
			return entity.Tokens{}, err
		} else {
			user = migratedUser
		}
	} else if err != nil {
		return entity.Tokens{}, failure.InvalidCredentials
	}

	if user.IsActivated == false {
		return entity.Tokens{}, failure.ProfileNotActivated
	}
	if user.IsBlocked == true {
		return entity.Tokens{}, failure.ProfileIsBlocked
	}

	if err = s.hashManager.ValidateByHash(password, user.Password); err != nil {
		return entity.Tokens{}, failure.InvalidCredentials
	}

	tokens, session, err := s.createSessionModel(user.Id, ip)
	if err != nil {
		return entity.Tokens{}, err
	}

	if err = s.repo.CreateSession(session); err != nil {
		return entity.Tokens{}, err
	}

	_ = s.repo.DeleteOldSessions(user.Id, maxSessionsCount)

	return tokens, nil
}

func (s *AuthService) SignOut(refreshToken string) error {
	return s.repo.DeleteSession(refreshToken)
}

func (s *AuthService) RefreshSession(refreshToken, ip string) (entity.Tokens, error) {
	user, err := s.repo.GetUserByRefreshToken(refreshToken)
	if err != nil {
		return entity.Tokens{}, err
	}

	if user.IsBlocked == true {
		_ = s.repo.DeleteSession(refreshToken)
		return entity.Tokens{}, failure.ProfileIsBlocked
	}

	tokens, session, err := s.createSessionModel(user.Id, ip)
	if err != nil {
		return entity.Tokens{}, err
	}

	return tokens, s.repo.UpdateSession(session, refreshToken)
}

func (s *AuthService) sendActivationLink(email string, activationLink uuid.UUID) error {
	if err := s.mailService.SendVerificationEmail(entity.VerificationEmailInput{
		Email:            email,
		Domain:           s.domain,
		VerificationCode: activationLink,
	}); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) migrateFromFirebase(credentials entity.Credentials) (entity.Profile, error) {
	firebaseUser, err := s.firebaseService.SignIn(credentials)
	if err != nil {
		return entity.Profile{}, failure.InvalidCredentials
	}

	credentials.Password, err = s.hashManager.Hash(credentials.Password)
	if err != nil {
		return entity.Profile{}, failure.Unknown
	}

	err = s.firebaseService.MigrateFromFirebase(credentials, firebaseUser)
	if err != nil {
		return entity.Profile{}, failure.UnableImportFirebaseProfile
	}

	user, err := s.repo.GetUserByEmail(credentials.Email)
	if err != nil {
		return entity.Profile{}, failure.InvalidCredentials
	}

	return user, nil
}

func (s *AuthService) createSessionModel(userId int, ip string) (entity.Tokens, entity.Session, error) {
	var (
		res entity.Tokens
		err error
	)
	res.AccessToken, err = s.tokenManager.NewJWT(strconv.Itoa(userId), s.accessTokenTTL)
	if err != nil {
		return entity.Tokens{}, entity.Session{}, failure.Unknown
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return entity.Tokens{}, entity.Session{}, failure.Unknown
	}

	return res, entity.Session{
		UserId:       userId,
		RefreshToken: res.RefreshToken,
		Ip:           ip,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}, nil
}

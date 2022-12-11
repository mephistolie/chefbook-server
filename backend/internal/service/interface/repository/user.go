package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"time"
)

type Auth interface {
	CreateUser(credentials entity.Credentials, activationLink uuid.UUID) (string, error)
	GetUserById(userId string) (entity.Profile, error)
	GetUserByEmail(email string) (entity.Profile, error)
	GetUserByRefreshToken(refreshToken string) (entity.Profile, error)
	GetUserActivationLink(userId string) (uuid.UUID, error)
	ActivateProfile(activationLink uuid.UUID) error
	ChangePassword(userId string, password string) error
	CreateSession(session entity.Session) error
	UpdateSession(session entity.Session, oldRefreshToken string) error
	DeleteSession(refreshToken string) error
	DeleteOldSessions(userId string, sessionsThreshold int) error
}

type Profile interface {
	SetUsername(userId string, username *string) error
	SetAvatarLink(userId string, url *string) error
	SetPremiumDate(userId string, expiresAt time.Time) error
	SetProfileCreationDate(userId string, creationTimestamp time.Time) error
	IncreaseBroccoins(userId string, broccoins int) error
}

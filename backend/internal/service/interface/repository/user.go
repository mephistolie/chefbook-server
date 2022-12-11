package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"time"
)

type Auth interface {
	CreateUser(credentials entity.Credentials, activationLink uuid.UUID) (uuid.UUID, error)
	GetUserById(userId uuid.UUID) (entity.Profile, error)
	GetUserByEmail(email string) (entity.Profile, error)
	GetUserByRefreshToken(refreshToken string) (entity.Profile, error)
	GetUserActivationLink(userId uuid.UUID) (uuid.UUID, error)
	ActivateProfile(activationLink uuid.UUID) error
	ChangePassword(userId uuid.UUID, password string) error
	CreateSession(session entity.Session) error
	UpdateSession(session entity.Session, oldRefreshToken string) error
	DeleteSession(refreshToken string) error
	DeleteOldSessions(userId uuid.UUID, sessionsThreshold int) error
}

type Profile interface {
	SetUsername(userId uuid.UUID, username *string) error
	SetAvatarLink(userId uuid.UUID, url *string) error
	SetPremiumDate(userId uuid.UUID, expiresAt time.Time) error
	SetProfileCreationDate(userId uuid.UUID, creationTimestamp time.Time) error
	IncreaseBroccoins(userId uuid.UUID, broccoins int) error
}

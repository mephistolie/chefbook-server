package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"time"
)

type Auth interface {
	CreateUser(credentials entity.Credentials, activationLink uuid.UUID) (int, error)
	GetUserById(userId int) (entity.Profile, error)
	GetUserByEmail(email string) (entity.Profile, error)
	GetUserByRefreshToken(refreshToken string) (entity.Profile, error)
	GetUserActivationLink(userId int) (uuid.UUID, error)
	ActivateProfile(activationLink uuid.UUID) error
	ChangePassword(userId int, password string) error
	CreateSession(session entity.Session) error
	UpdateSession(session entity.Session, oldRefreshToken string) error
	DeleteSession(refreshToken string) error
	DeleteOldSessions(userId, sessionsThreshold int) error
}

type Profile interface {
	SetUsername(userId int, username *string) error
	SetAvatarLink(userId int, url *string) error
	SetPremiumDate(userId int, expiresAt time.Time) error
	SetProfileCreationDate(userId int, creationTimestamp time.Time) error
	IncreaseBroccoins(userId, broccoins int) error
}

package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Auth interface {
	SignUp(credentials entity.Credentials) (int, error)
	ActivateProfile(activationLink uuid.UUID) error
	SignIn(credentials entity.Credentials, ip string) (entity.Tokens, error)
	SignOut(refreshToken string) error
	RefreshSession(refreshToken, ip string) (entity.Tokens, error)
}

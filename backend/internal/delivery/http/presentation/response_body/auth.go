package response_body

import (
	"github.com/mephistolie/chefbook-server/internal/entity"
	"time"
)

type Tokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func NewTokens(tokens entity.Tokens) Tokens {
	return Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}
}

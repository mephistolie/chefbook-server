package response_body

import "github.com/mephistolie/chefbook-server/internal/entity"

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokens(tokens entity.Tokens) Tokens {
	return Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

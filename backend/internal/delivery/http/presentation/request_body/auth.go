package request_body

import (
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"unicode"
)

type Credentials struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

func (c *Credentials) Validate() error {
	return validatePassword(c.Password)
}

func validatePassword(password string) error  {
	lower := false
	upper := false
	number := false
	for _, c := range password {
		switch {
		case unicode.IsLower(c):
			lower = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsNumber(c):
			number = true
		case c == ' ':
			return failure.InvalidBody
		default:
		}
	}
	if !lower || !upper || !number {
		return failure.InvalidBody
	}
	return nil
}

func (c *Credentials) Entity() entity.Credentials {
	return entity.Credentials{
		Email:    c.Email,
		Password: c.Password,
	}
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
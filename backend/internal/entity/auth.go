package entity

import (
	"github.com/google/uuid"
	"time"
)

type Credentials struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Session struct {
	UserId       int
	RefreshToken string
	Ip           string
	ExpiresAt    time.Time
}

type VerificationEmailInput struct {
	Email            string
	Name             string
	VerificationCode uuid.UUID
	Domain           string
}

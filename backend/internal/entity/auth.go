package entity

import (
	"github.com/google/uuid"
	"time"
)

type Credentials struct {
	Id       *uuid.UUID
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type Session struct {
	UserId       uuid.UUID
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

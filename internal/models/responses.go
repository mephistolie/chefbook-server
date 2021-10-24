package models

import "errors"

var (
	RespActivationLink   = "profile activation link sent to email"
	RespProfileActivated = "profile activated"

	ErrUserAlreadyExists   = errors.New("user with such email already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrProfileNotActivated = errors.New("profile not activated. check your email")
	ErrAuthentication      = errors.New("invalid sign in data")
)

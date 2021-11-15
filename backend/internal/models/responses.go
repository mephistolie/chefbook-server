package models

import "errors"

var (
	RespActivationLink   = "profile activation link sent to email"
	RespProfileActivated = "profile activated"
	RespSignOutSuccessfully = "sign out successfully"

	RespRecipeAdded = "recipe added"
	RespRecipeUpdated = "recipe updated"
	RespRecipeDeleted = "recipe deleted"

	ErrInvalidInput = errors.New("invalid input")
	ErrAccessDenied   = errors.New("recipe access denied")

	ErrUserAlreadyExists   = errors.New("user with such email already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrProfileNotActivated = errors.New("profile not activated. check your email")
	ErrProfileIsBlocked    = errors.New("profile is blocked")
	ErrAuthentication      = errors.New("invalid sign in data")

	ErrEmptyAuthHeader   = errors.New("empty auth header")
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrEmptyToken        = errors.New("token is empty")
	ErrSessionExpired    = errors.New("session expired")
	ErrSessionNotFound   = errors.New("session not found")

	ErrUserIdNotFound = errors.New("user id not found")
	ErrInvalidUserId  = errors.New("invalid user id")

	ErrInvalidRecipeInput = errors.New("invalid recipe input")
	ErrRecipeNotFound = errors.New("recipe not found")
)

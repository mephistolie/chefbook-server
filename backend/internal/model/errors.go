package model

import "errors"

var(
	ErrInvalidInput         = errors.New("invalid input")
	ErrInvalidFileInput     = errors.New("invalid input. File must be < 500KB")
	ErrFileTypeNotSupported = errors.New("file type isn't supported")
	ErrAccessDenied         = errors.New("recipe access denied")

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

	ErrUnableDeleteAvatar  = errors.New("unable to delete avatar")
	ErrNoKey               = errors.New("encrypted key not found")
	ErrUnableDeleteUserKey = errors.New("unable to delete user key")

	ErrNotOwner                  = errors.New("you aren't owner of this recipe")
	ErrInvalidRecipeInput        = errors.New("invalid recipe input")
	ErrRecipeNotFound            = errors.New("recipe not found")
	ErrUnableAddRecipe           = errors.New("unable to add recipe to recipe book")
	ErrRecipeNotInRecipeBook     = errors.New("recipe isn't in recipe book")
	ErrRecipeLikeSetAlready      = errors.New("recipe like status already set")
	ErrUnableDeleteRecipePicture = errors.New("unable to delete picture")
	ErrUnableDeleteRecipeKey     = errors.New("unable to delete recipe key")
	ErrUnableGetRandomRecipe     = errors.New("unable to found random recipe with request parameters")

	ErrCategoryNotFound = errors.New("category not found")
)

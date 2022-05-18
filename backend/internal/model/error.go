package model

import "errors"

var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrBigFile              = errors.New("invalid input. File must be < 500KB")
	ErrFileTypeNotSupported = errors.New("file type isn't supported")
	ErrUnableUploadFile     = errors.New("unable to upload file")
	ErrAccessDenied         = errors.New("access denied")

	ErrUnableSendEmail       = errors.New("unable to send email")
	ErrUserAlreadyExists     = errors.New("user with such email already exists")
	ErrProfileNotActivated   = errors.New("profile not activated. check your email")
	ErrInvalidActivationLink = errors.New("invalid activation link")
	ErrProfileIsBlocked      = errors.New("profile is blocked")
	ErrInvalidAuthData       = errors.New("invalid auth data")

	ErrFirebaseImport = errors.New("can't import old profile")

	ErrEmptyAuthHeader   = errors.New("empty auth header")
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrEmptyToken        = errors.New("token is empty")
	ErrInvalidToken      = errors.New("token is invalid")
	ErrSessionExpired    = errors.New("session expired")
	ErrSessionNotFound   = errors.New("session not found")

	ErrUserNotFound   = errors.New("user not found")
	ErrUserIdNotFound = errors.New("user id not found")
	ErrInvalidUserId  = errors.New("invalid user id")

	ErrUnableSetAvatar    = errors.New("unable set avatar")
	ErrUnableDeleteAvatar = errors.New("unable to delete avatar")

	ErrNoKey               = errors.New("encrypted key not found")
	ErrUnableSetUserKey    = errors.New("unable to set user key")
	ErrUnableDeleteUserKey = errors.New("unable to delete user key")

	ErrNotOwner                  = errors.New("you aren't owner of this recipe")
	ErrRecipeNotFound            = errors.New("recipe not found")
	ErrUnableAddRecipe           = errors.New("unable to add recipe to recipe book")
	ErrRecipeNotInRecipeBook     = errors.New("recipe isn't in recipe book")
	ErrUnableDeleteRecipePicture = errors.New("unable to delete picture")
	ErrUnableDeleteRecipeKey     = errors.New("unable to delete recipe key")
	ErrUnableGetRandomRecipe     = errors.New("unable to found random recipe with request parameters")

	ErrCategoryNotFound = errors.New("category not found")

	ErrShoppingListNotFound = errors.New("shopping list not found")
)

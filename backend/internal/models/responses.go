package models

import "errors"

var (
	RespActivationLink   = "profile activation link has been sent to email"
	RespProfileActivated = "profile is activated"
	RespSignOutSuccessfully = "signed out successfully"

	RespUsernameChanged = "username successfully changed"
	RespAvatarDeleted = "avatar has been deleted"

	RespRecipeAdded = "recipe has been added"
	RespRecipeUpdated = "recipe has been updated"
	RespRecipeDeleted = "recipe has been deleted"
	RespCategoriesUpdated = "categories has been updated"
	RespFavouriteStatusUpdated = "favourite status has been updated"
	RespRecipeLikeSet = "recipe like status has been set"

	RespCategoryAdded = "category has been added"
	RespCategoryUpdated = "category has been updated"
	RespCategoryDeleted = "category has been deleted"

	RespShoppingListUpdated = "shopping list has been updated"

	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidFileInput = errors.New("invalid input. File must be < 500KB")
	ErrFileTypeNotSupported = errors.New("file type isn't supported")
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

	ErrUnableDeleteAvatar = errors.New("unable to delete avatar")

	ErrInvalidRecipeInput = errors.New("invalid recipe input")
	ErrRecipeNotFound = errors.New("recipe not found")
	ErrRecipeNotInRecipeBook = errors.New("recipe isn't in recipe book")
	ErrRecipeLikeSetAlready = errors.New("recipe like status already set")

	ErrCategoryNotFound = errors.New("category not found")
)

package response_body

import (
	"chefbook-server/internal/entity/failure"
)

const (
	errTypeAccessDenied        = "access_denied"
	errTypeInvalidAccessToken  = "unauthorized"
	errTypeInvalidRefreshToken = "invalid_refresh_token"

	errTypeInvalidBody = "invalid_body"
	errTypeBigFile     = "big_file"
	errTypeNotFound    = "not_found"

	errTypeInvalidCredentials    = "invalid_credentials"
	errTypeProfileNotActivated   = "profile_not_activated"
	errTypeInvalidActivationLink = "invalid_activation_link"
	errTypeUserExists            = "user_exists"
	errTypeUserBlocked           = "user_blocked"

	errFirebaseProfileImport = "import_old_profile_failed"

	errInvalidRecipe   = "invalid_recipe"
	errNotInRecipeBook = "not_in_recipe_book"

	errTypeUnknown = "unknown_error"
)

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewError(err error) Error {
	errType := errTypeUnknown
	switch err {
	case failure.AccessDenied, failure.NotOwner:
		errType = errTypeAccessDenied
	case failure.EmptyAuthHeader, failure.InvalidAuthHeader, failure.EmptyToken, failure.InvalidToken,
		failure.SessionExpired:
		errType = errTypeInvalidAccessToken
	case failure.UserNotFound, failure.RecipeNotFound, failure.CategoryNotFound, failure.ActivationLinkNotFound,
		failure.NoKey, failure.ShoppingListNotFound, failure.UnableGetRandomRecipe:
		errType = errTypeNotFound
	case failure.SessionNotFound:
		errType = errTypeInvalidRefreshToken
	case failure.InvalidBody, failure.UnsupportedFileType, failure.EmptyRecipeName, failure.EmptyIngredients, failure.EmptyCooking,
		failure.InvalidUserId, failure.TooLongRecipeName, failure.TooLongRecipeDescription, failure.TooLongIngredientItemText,
		failure.InvalidIngredientItemType, failure.InvalidCookingItemType, failure.InvalidEncryptionType:
		errType = errTypeInvalidBody
	case failure.InvalidFileSize:
		errType = errTypeBigFile
	case failure.InvalidCredentials:
		errType = errTypeInvalidCredentials
	case failure.ProfileNotActivated:
		errType = errTypeProfileNotActivated
	case failure.InvalidActivationLink:
		errType = errTypeInvalidActivationLink
	case failure.UserAlreadyExists:
		errType = errTypeUserExists
	case failure.ProfileIsBlocked:
		errType = errTypeUserBlocked
	case failure.UnableImportFirebaseProfile:
		errType = errFirebaseProfileImport
	case failure.InvalidRecipe:
		errType = errInvalidRecipe
	case failure.RecipeNotInRecipeBook:
		errType = errNotInRecipeBook
	}

	return Error{
		Error:   errType,
		Message: err.Error(),
	}
}

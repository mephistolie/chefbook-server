package response_body

import (
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"net/http"
)

const (
	errTypeAccessDenied        = "access_denied"
	errTypeUnauthorized        = "unauthorized"
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

func NewError(err error) (int, Error) {
	statusCode := http.StatusBadRequest
	errType := errTypeUnknown
	switch err {
	case failure.EmptyAuthHeader, failure.InvalidAuthHeader, failure.EmptyToken, failure.InvalidToken,
		failure.SessionExpired:
		statusCode = http.StatusUnauthorized
		errType = errTypeUnauthorized
	case failure.SessionNotFound:
		statusCode = http.StatusUnauthorized
		errType = errTypeInvalidRefreshToken
	case failure.AccessDenied, failure.NotOwner:
		statusCode = http.StatusForbidden
		errType = errTypeAccessDenied
	case failure.UserNotFound, failure.RecipeNotFound, failure.CategoryNotFound, failure.ActivationLinkNotFound,
		failure.NoKey, failure.ShoppingListNotFound, failure.UnableGetRandomRecipe:
		statusCode = http.StatusNotFound
		errType = errTypeNotFound
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

	return statusCode, Error{
		Error:   errType,
		Message: err.Error(),
	}
}

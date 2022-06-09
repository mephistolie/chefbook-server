package response_body

import (
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
)

const (
	errTypeAccessDenied        = "ACCESS_DENIED"
	errTypeInvalidAccessToken  = "INVALID_ACCESS_TOKEN"
	errTypeInvalidRefreshToken = "INVALID_REFRESH_TOKEN"

	errTypeInvalidBody = "INVALID_BODY"
	errTypeBigFile     = "BIG_FILE"
	errTypeNotFound    = "NOT_FOUND"

	errTypeUnableSendMail        = "UNABLE_SEND_MAIL"
	errTypeInvalidCredentials    = "INVALID_CREDENTIALS"
	errTypeProfileNotActivated   = "PROFILE_NOT_ACTIVATED"
	errTypeInvalidActivationLink = "INVALID_ACTIVATION_LINK"
	errTypeUserExists            = "USER_EXISTS"
	errTypeUserBlocked           = "USER_BLOCKED"

	errFirebaseProfileImport = "IMPORT_OLD_PROFILE_FAILED"

	errInvalidRecipe = "INVALID_RECIPE"
	errNotInRecipeBook = "NOT_IN_RECIPE_BOOK"

	errTypeUnknown = "UNKNOWN_ERROR"
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
	case failure.UnableSendEmail:
		errType = errTypeUnableSendMail
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

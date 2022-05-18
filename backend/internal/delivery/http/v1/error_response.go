package v1

import "github.com/mephistolie/chefbook-server/internal/model"

const (
	errTypeAccessDenied        = "ACCESS_DENIED"
	errTypeInvalidAccessToken  = "INVALID_ACCESS_TOKEN"
	errTypeInvalidRefreshToken = "INVALID_REFRESH_TOKEN"

	errTypeInvalidInput = "INVALID_INPUT"
	errTypeBigFile      = "BIG_FILE"
	errTypeNotFound     = "NOT_FOUND"

	errTypeUnableSendMail        = "UNABLE_SEND_MAIL"
	errTypeInvalidAuthData       = "INVALID_AUTH_DATA"
	errTypeProfileNotActivated   = "PROFILE_NOT_ACTIVATED"
	errTypeInvalidActivationLink = "INVALID_ACTIVATION_LINK"
	errTypeUserExists            = "USER_EXISTS"
	errTypeUserBlocked           = "USER_BLOCKED"

	errFirebaseProfileImport = "IMPORT_OLD_PROFILE_FAILED"

	errNotInRecipeBook = "NOT_IN_RECIPE_BOOK"

	errTypeUnknown = "UNKNOWN_ERROR"
)

func getErrorResponseBody(err error) ErrorResponse {
	errType := errTypeUnknown
	switch err {
	case model.ErrAccessDenied, model.ErrNotOwner:
		errType = errTypeAccessDenied
	case model.ErrEmptyAuthHeader, model.ErrInvalidAuthHeader, model.ErrEmptyToken, model.ErrInvalidToken,
		model.ErrSessionExpired:
		errType = errTypeInvalidAccessToken
	case model.ErrUserNotFound, model.ErrUserIdNotFound, model.ErrRecipeNotFound, model.ErrCategoryNotFound,
		model.ErrNoKey, model.ErrShoppingListNotFound, model.ErrUnableGetRandomRecipe:
		errType = errTypeNotFound
	case model.ErrSessionNotFound:
		errType = errTypeInvalidRefreshToken
	case model.ErrInvalidInput, model.ErrFileTypeNotSupported, model.ErrInvalidUserId:
		errType = errTypeInvalidInput
	case model.ErrBigFile:
		errType = errTypeBigFile
	case model.ErrUnableSendEmail:
		errType = errTypeUnableSendMail
	case model.ErrInvalidAuthData:
		errType = errTypeInvalidAuthData
	case model.ErrProfileNotActivated:
		errType = errTypeProfileNotActivated
	case model.ErrInvalidActivationLink:
		errType = errTypeInvalidActivationLink
	case model.ErrUserAlreadyExists:
		errType = errTypeUserExists
	case model.ErrProfileIsBlocked:
		errType = errTypeUserBlocked
	case model.ErrFirebaseImport:
		errType = errFirebaseProfileImport
	case model.ErrRecipeNotInRecipeBook:
		errType = errNotInRecipeBook
	}

	return ErrorResponse{
		Error:   errType,
		Message: err.Error(),
	}
}

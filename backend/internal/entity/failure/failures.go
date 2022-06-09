package failure

import "errors"

var (
	Unknown = errors.New("unknown error")

	UnableCreateProfile = errors.New("unable to create profile")
	UnableDeleteSession = errors.New("unable to delete session")

	InvalidBody         = errors.New("invalid request body")
	InvalidFileSize     = errors.New("invalid file size")
	UnsupportedFileType = errors.New("unsupported file type")
	UnableUploadFile    = errors.New("unable to upload file")
	UnableDeleteFile    = errors.New("unable delete file")
	AccessDenied        = errors.New("access denied")

	UnableSendEmail       = errors.New("unable to send email")
	UserAlreadyExists     = errors.New("user with such email already exists")
	ProfileNotActivated   = errors.New("profile not activated. check your email")
	InvalidActivationLink = errors.New("invalid activation link")
	ProfileIsBlocked      = errors.New("profile is blocked")
	InvalidCredentials    = errors.New("invalid credentials")

	UnableImportFirebaseProfile = errors.New("can't import old profile")

	EmptyAuthHeader   = errors.New("empty auth header")
	InvalidAuthHeader = errors.New("invalid auth header")
	EmptyToken        = errors.New("token is empty")
	InvalidToken      = errors.New("token is invalid")
	SessionExpired    = errors.New("session expired")
	SessionNotFound   = errors.New("session not found")

	UserNotFound           = errors.New("user not found")
	InvalidUserId          = errors.New("invalid user id")
	ActivationLinkNotFound = errors.New("activation link not found")

	UnableSetAvatar = errors.New("unable set avatar")

	NoKey = errors.New("encrypted key not found")

	EmptyRecipeName           = errors.New("empty recipe name")
	EmptyIngredients          = errors.New("no ingredients")
	EmptyCooking              = errors.New("no cooking")
	TooLongRecipeName         = errors.New("too long recipe name; maximum is 100 symbols")
	TooLongRecipeDescription  = errors.New("too long recipe description; maximum is 1500 symbols")
	TooLongIngredientItemText = errors.New("too long decrypted ingredient item text; maximum is 100 symbols")
	InvalidIngredientItemType = errors.New("invalid ingredient type")
	InvalidCookingItemType    = errors.New("invalid cooking step type")
	InvalidEncryptionType     = errors.New("recipe input doesn't match its encryption state")

	NotOwner              = errors.New("you aren't owner of this recipe")
	UnableCreateRecipe    = errors.New("unable to create recipe")
	RecipeNotFound        = errors.New("recipe not found")
	InvalidRecipe         = errors.New("recipe has invalid fields")
	UnableAddRecipe       = errors.New("unable to add recipe to recipe book")
	RecipeNotInRecipeBook = errors.New("recipe isn't in recipe book")
	UnableGetRandomRecipe = errors.New("unable to found random recipe with request parameters")

	UnableAddCategory = errors.New("unable to add category")
	CategoryNotFound  = errors.New("category not found")

	ShoppingListNotFound = errors.New("shopping list not found")
)

package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/model"
	"github.com/mephistolie/chefbook-server/internal/repository/postgres"
	"github.com/mephistolie/chefbook-server/internal/repository/s3"
	"github.com/minio/minio-go/v7"
	"time"
)

type Auth interface {
	CreateUser(user model.AuthData, activationLink uuid.UUID) (int, error)
	GetUserById(id int) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
	GetUserByCredentials(email, password string) (model.User, error)
	GetByRefreshToken(refreshToken string) (model.User, error)
	ActivateUser(activationLink uuid.UUID) error
	CreateSession(session model.Session) error
	UpdateSession(session model.Session, oldRefreshToken string) error
	DeleteSession(refreshToken string) error
	ChangePassword(user model.AuthData) error
}

type Profile interface {
	SetUsername(userId int, username string) error
	SetAvatar(userId int, url string) error
	SetPremiumDate(userId int, expiresAt time.Time) error
	SetProfileCreationDate(userId int, creationTimestamp time.Time) error
	IncreaseBroccoins(userId, broccoins int) error
	ReduceBroccoins(userId, broccoins int) error
}

type RecipeCrud interface {
	GetRecipesInfoByRequest(params model.RecipesRequestParams) ([]model.RecipeInfo, error)
	GetRecipeOwnerId(recipeId int) (int, error)
	CreateRecipe(recipe model.Recipe) (int, error)
	AddRecipeToRecipeBook(recipeId, userId int) error
	GetRecipe(recipeId int) (model.Recipe, error)
	GetRecipeWithUserFields(recipeId int, userId int) (model.Recipe, error)
	UpdateRecipe(recipe model.Recipe) error
	DeleteRecipe(recipeId int) error
	DeleteRecipeFromRecipeBook(recipeId, userId int) error
}

type RecipeInteraction interface {
	SetRecipeCategories(categoriesIds []int, recipeId, userId int) error
	SetRecipeFavourite(recipeId, userId int, isFavourite bool) error
	SetRecipeLiked(recipeId, userId int, isLiked bool) error
}

type RecipeSharing interface {
	GetRecipeUserList(recipeId int) ([]model.UserInfo, error)
	SetUserPublicKeyForRecipe(recipeId int, userId int, userKey string) error
	SetUserPrivateKeyForRecipe(recipeId int, userId int, userKey string) error
}

type Encryption interface {
	GetUserKey(userId int) (string, error)
	SetUserKey(userId int, url string) error
	GetRecipeKey(recipeId int) (string, error)
	SetRecipeKey(recipeId int, url string) error
}

type Categories interface {
	GetUserCategories(userId int) ([]model.Category, error)
	AddCategory(category model.Category) (int, error)
	GetCategoryById(categoryId int) (model.Category, error)
	UpdateCategory(category model.Category) error
	DeleteCategory(categoryId, userId int) error
	GetRecipeCategories(recipeId, userId int) ([]int, error)
}

type ShoppingList interface {
	GetShoppingList(userId int) (model.ShoppingList, error)
	SetShoppingList(shoppingList model.ShoppingList, userId int) error
}

type Files interface {
	UploadAvatar(ctx context.Context, userId int, input model.MultipartFileInfo) (string, error)
	UploadUserKey(ctx context.Context, userId int, input model.MultipartFileInfo) (string, error)
	GetRecipePictures(ctx context.Context, recipeId int) []string
	UploadRecipePicture(ctx context.Context, recipeId int, input model.MultipartFileInfo) (string, error)
	UploadRecipeKey(ctx context.Context, recipeId int, input model.MultipartFileInfo) (string, error)
	GetRecipePictureLink(recipeId int, pictureName string) string
	GetRecipeKeysLink(recipeId int, pictureName string) string
	DeleteFile(ctx context.Context, url string) error
}

type Repository struct {
	Auth
	Profile
	RecipeCrud
	RecipeInteraction
	RecipeSharing
	Encryption
	Categories
	ShoppingList
	Files
}

func NewRepositories(db *sqlx.DB, client *minio.Client) *Repository {
	return &Repository{
		Auth:              postgres.NewUsersPostgres(db),
		Profile:           postgres.NewProfilePostgres(db),
		RecipeCrud:        postgres.NewRecipesPostgres(db),
		RecipeInteraction: postgres.NewRecipeInteractionPostgres(db),
		RecipeSharing:     postgres.NewRecipeSharingPostgres(db),
		Encryption:        postgres.NewEncryptionPostgres(db),
		Categories:        postgres.NewCategoriesPostgres(db),
		ShoppingList:      postgres.NewShoppingListPostgres(db),
		Files:             s3.NewAWSFileManager(client),
	}
}

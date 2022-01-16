package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository/postgres"
	"github.com/mephistolie/chefbook-server/internal/repository/s3"
	"github.com/minio/minio-go/v7"
	"time"
)

type Users interface {
	CreateUser(user models.AuthData, activationLink uuid.UUID) (int, error)
	GetUserById(id int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByCredentials(email, password string) (models.User, error)
	GetByRefreshToken(refreshToken string) (models.User, error)
	ActivateUser(activationLink uuid.UUID) error
	CreateSession(session models.Session) error
	UpdateSession(session models.Session, oldRefreshToken string) error
	DeleteSession(refreshToken string) error
	ChangePassword(user models.AuthData) error

	SetUserName(userId int, username string) error
	SetUserAvatar(userId int, url string) error
	SetPremiumDate(userId int, expiresAt time.Time) error

	GetUserKey(userId int) (string, error)
	SetUserKey(userId int, url string) error
}

type Recipes interface {
	GetRecipesByUser(userId int) ([]models.Recipe, error)
	GetRecipeOwnerId(recipeId int) (int, error)
	CreateRecipe(recipe models.Recipe) (int, error)
	GetRecipeById(recipeId int, userId int) (models.Recipe, error)
	UpdateRecipe(recipe models.Recipe, userId int) error
	DeleteRecipe(recipeId int) error
	DeleteRecipeLink(recipeId, userId int) error
	SetRecipeCategories(categoriesIds []int, recipeId, userId int) error
	MarkRecipeFavourite(recipeId, userId int, isFavourite bool) error
	SetRecipeLike(recipeId, userId int, isLiked bool) error
	SetRecipePreview(recipeId int, url string)  error
	GetRecipeKey(recipeId int) (string, error)
    SetRecipeKey(recipeId int, url string) error
}

type Categories interface {
	GetCategoriesByUser(userId int) ([]models.Category, error)
	AddCategory(category models.Category) (int, error)
	GetCategoryById(categoryId, userId int) (models.Category, error)
	UpdateCategory(category models.Category) error
	DeleteCategory(categoryId, userId int) error
	GetRecipeCategories(recipeId, userId int) ([]int, error)
}

type ShoppingList interface {
	GetShoppingList(userId int) (models.ShoppingList, error)
	SetShoppingList(shoppingList models.ShoppingList, userId int) error
}

type Files interface {
	UploadAvatar(ctx context.Context, userId int, input s3.UploadInput) (string, error)
	UploadUserKey(ctx context.Context, userId int, input s3.UploadInput) (string, error)
	UploadRecipePicture(ctx context.Context, recipeId int, input s3.UploadInput) (string, error)
	UploadRecipeKey(ctx context.Context, recipeId int, input s3.UploadInput) (string, error)
	GetRecipePictureLink(recipeId int, pictureName string) string
	GetRecipeKeysLink(recipeId int, pictureName string) string
	DeleteFile(ctx context.Context, url string) error
}

type Repository struct {
	Users
	Recipes
	Categories
	ShoppingList
	Files
}

func NewRepositories(db *sqlx.DB, client *minio.Client) *Repository {
	return &Repository{
		Users:        postgres.NewUsersPostgres(db),
		Recipes:      postgres.NewRecipesPostgres(db),
		Categories:   postgres.NewCategoriesPostgres(db),
		ShoppingList: postgres.NewShoppingListPostgres(db),
		Files:        s3.NewAWSFileManager(client),
	}
}

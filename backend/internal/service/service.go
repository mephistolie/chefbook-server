package service

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/config"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"github.com/mephistolie/chefbook-server/pkg/auth"
	"github.com/mephistolie/chefbook-server/pkg/cache"
	"github.com/mephistolie/chefbook-server/pkg/hash"
	"github.com/mephistolie/chefbook-server/pkg/mail"
	"time"
)

type Auth interface {
	SignUp(authInput models.AuthData) (int, error)
	ActivateUser(activationLink uuid.UUID) error
	SignIn(authInput models.AuthData, ip string) (models.Tokens, error)
	SignOut(refreshToken string) error
	RefreshSession(refreshToken, ip string) (models.Tokens, error)
}

type Firebase interface {
	FirebaseSignIn(authData models.AuthData) (models.FirebaseUser, error)
}

type Users interface {
	GetUserInfo(userId int) (models.User, error)
	SetUserName(userId int, username string) error
	UploadAvatar(ctx context.Context, userId int, file *bytes.Reader, size int64, contentType string) (string, error)
	DeleteAvatar(ctx context.Context, userId int) error

	GetUserKey(userId int) (string, error)
	UploadUserKey(ctx context.Context, userId int, file *bytes.Reader, size int64, contentType string) (string, error)
	DeleteUserKey(ctx context.Context, userId int) error
}

type VerificationEmailInput struct {
	Email            string
	Name             string
	VerificationCode uuid.UUID
	Domain           string
}

type Mails interface {
	SendVerificationEmail(input VerificationEmailInput) error
}

type Recipes interface {
	GetRecipesByUser(userId int) ([]models.Recipe, error)
	AddRecipe(recipe models.Recipe) (int, error)
	GetRecipeById(recipeId, userId int) (models.Recipe, error)
	UpdateRecipe(recipe models.Recipe, userId int) error
	DeleteRecipe(recipeId, userId int) error
	SetRecipeCategories(input models.RecipeCategoriesInput) error
	MarkRecipeFavourite(input models.FavouriteRecipeInput) error
	SetRecipeLike(input models.RecipeLikeInput) error
	UploadRecipePicture(ctx context.Context, recipeId, userId int, file *bytes.Reader, size int64, contentType string) (string, error)
	DeleteRecipePicture(ctx context.Context, recipeId, userId int, pictureName string) error

	GetRecipeKey(recipeId, userId int) (string, error)
	UploadRecipeKey(ctx context.Context, recipeId, userId int, file *bytes.Reader, size int64, contentType string) (string, error)
	DeleteRecipeKey(ctx context.Context, recipeId, userId int) error
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
	AddToShoppingList(newPurchases []models.Purchase, userId int) error
}

type Service struct {
	Auth
	Firebase
	Users
	Mails
	Recipes
	Categories
	ShoppingList
}

type Dependencies struct {
	Repos           *repository.Repository
	Cache           cache.Cache
	HashManager     hash.HashManager
	TokenManager    auth.TokenManager
	MailSender      mail.Sender
	MailConfig      config.MailConfig
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	CacheTTL        int64
	Environment     string
	Domain          string
	FirebaseApiKey  string
	FirestoreClient firestore.Client
}

func NewServices(dependencies Dependencies) *Service {

	mailService := NewMailService(dependencies.MailSender, dependencies.MailConfig, dependencies.Cache)
	firebaseService := NewFirebaseService(dependencies.FirebaseApiKey, dependencies.Repos.Users, dependencies.Repos.Recipes, dependencies.Repos.Categories, dependencies.Repos.ShoppingList, dependencies.FirestoreClient)

	return &Service{
		Auth: NewAuthService(dependencies.Repos, *firebaseService, dependencies.HashManager, dependencies.TokenManager,
			dependencies.AccessTokenTTL, dependencies.RefreshTokenTTL, mailService, dependencies.Domain),
		Firebase:     firebaseService,
		Users:        NewUsersService(dependencies.Repos.Users, dependencies.Repos.Files),
		Mails:        mailService,
		Recipes:      NewRecipesService(dependencies.Repos.Recipes, dependencies.Repos.Categories, dependencies.Repos.Files),
		Categories:   NewCategoriesService(dependencies.Repos),
		ShoppingList: NewShoppingListService(dependencies.Repos),
	}
}

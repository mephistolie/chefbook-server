package service

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/config"
	"github.com/mephistolie/chefbook-server/internal/model"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"github.com/mephistolie/chefbook-server/pkg/auth"
	"github.com/mephistolie/chefbook-server/pkg/cache"
	"github.com/mephistolie/chefbook-server/pkg/hash"
	"github.com/mephistolie/chefbook-server/pkg/mail"
	"time"
)

type Auth interface {
	SignUp(authInput model.AuthData) (int, error)
	ActivateUser(activationLink uuid.UUID) error
	SignIn(authInput model.AuthData, ip string) (model.Tokens, error)
	SignOut(refreshToken string) error
	RefreshSession(refreshToken, ip string) (model.Tokens, error)
}

type Firebase interface {
	SignIn(authData model.AuthData) (model.FirebaseUser, error)
}

type Profile interface {
	GetUserInfo(userId int) (model.User, error)
	SetUsername(userId int, username string) error
	UploadAvatar(ctx context.Context, userId int, file model.MultipartFileInfo) (string, error)
	DeleteAvatar(ctx context.Context, userId int) error
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

type RecipesCrud interface {
	GetRecipesInfoByRequest(params model.RecipesRequestParams) ([]model.RecipeInfo, error)
	CreateRecipe(recipe model.Recipe) (int, error)
	AddRecipeToRecipeBook(recipeId, userId int) error
	GetRecipeById(recipeId, userId int) (model.Recipe, error)
	UpdateRecipe(recipe model.Recipe) error
	DeleteRecipe(recipeId, userId int) error
}

type RecipeInteraction interface {
	SetRecipeCategories(input model.RecipeCategoriesInput) error
	SetRecipeFavourite(input model.FavouriteRecipeInput) error
	SetRecipeLiked(input model.RecipeLikeInput) error
}

type RecipeSharing interface {
	GetRecipeUserList(recipeId, userId int) ([]model.UserInfo, error)
	SetUserPublicKeyForRecipe(recipeId int, userId int, userKey string) error
	SetUserPrivateKeyForRecipe(recipeId int, userId int, userKey string) error
	DeleteUserAccessToRecipe(recipeId, userId, requesterId int) error
}

type RecipePictures interface {
	GetRecipePictures(ctx context.Context, recipeId int, userId int) ([]string, error)
	UploadRecipePicture(ctx context.Context, recipeId, userId int, file model.MultipartFileInfo) (string, error)
	DeleteRecipePicture(ctx context.Context, recipeId, userId int, pictureName string) error
}

type Encryption interface {
	GetUserKeyLink(userId int) (string, error)
	UploadUserKey(ctx context.Context, userId int, file model.MultipartFileInfo) (string, error)
	DeleteUserKey(ctx context.Context, userId int) error
	GetRecipeKey(recipeId, userId int) (string, error)
	UploadRecipeKey(ctx context.Context, recipeId, userId int, file model.MultipartFileInfo) (string, error)
	DeleteRecipeKey(ctx context.Context, recipeId, userId int) error
}

type Categories interface {
	GetUserCategories(userId int) ([]model.Category, error)
	GetRecipeCategories(recipeId, userId int) ([]int, error)
	AddCategory(category model.Category) (int, error)
	GetCategoryById(categoryId, userId int) (model.Category, error)
	UpdateCategory(category model.Category) error
	DeleteCategory(categoryId, userId int) error
}

type ShoppingList interface {
	GetShoppingList(userId int) (model.ShoppingList, error)
	SetShoppingList(shoppingList model.ShoppingList, userId int) error
	AddToShoppingList(newPurchases []model.Purchase, userId int) error
}

type Service struct {
	Auth
	Firebase
	Profile
	Mails
	RecipesCrud
	RecipeInteraction
	RecipeSharing
	RecipePictures
	Encryption
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
	FirebaseApp     firebase.App
}

func NewServices(dependencies Dependencies) *Service {

	mailService := NewMailService(dependencies.MailSender, dependencies.MailConfig, dependencies.Cache)
	firebaseService := NewFirebaseService(dependencies.FirebaseApiKey, dependencies.Repos.Auth, dependencies.Repos.Profile,
		dependencies.Repos.RecipeCrud, dependencies.Repos.RecipeInteraction, dependencies.Repos.Categories, dependencies.Repos.ShoppingList,
		dependencies.FirebaseApp)

	return &Service{
		Auth: NewAuthService(dependencies.Repos, *firebaseService, dependencies.HashManager, dependencies.TokenManager,
			dependencies.AccessTokenTTL, dependencies.RefreshTokenTTL, mailService, dependencies.Domain),
		Firebase:          firebaseService,
		Profile:           NewUsersService(dependencies.Repos.Auth, dependencies.Repos.Profile, dependencies.Repos.Files),
		Mails:             mailService,
		RecipesCrud:       NewRecipesService(dependencies.Repos.RecipeCrud, dependencies.Repos.Categories),
		RecipeInteraction: NewRecipeInteractionService(dependencies.Repos.RecipeInteraction),
		RecipeSharing:     NewRecipeSharingService(dependencies.Repos.RecipeCrud, dependencies.Repos.RecipeSharing),
		RecipePictures:    NewRecipePicturesService(dependencies.Repos.RecipeCrud, dependencies.Repos.Files),
		Encryption:        NewEncryptionService(dependencies.Repos.Encryption, dependencies.Repos.RecipeCrud, dependencies.Repos.Files),
		Categories:        NewCategoriesService(dependencies.Repos),
		ShoppingList:      NewShoppingListService(dependencies.Repos),
	}
}

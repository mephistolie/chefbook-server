package service

import (
	"github.com/mephistolie/chefbook-server/internal/app/dependencies/repository"
	"github.com/mephistolie/chefbook-server/internal/config"
	"github.com/mephistolie/chefbook-server/internal/service"
	"github.com/mephistolie/chefbook-server/pkg/auth"
	"github.com/mephistolie/chefbook-server/pkg/hash"
	"github.com/mephistolie/chefbook-server/pkg/mail"
	"time"
)

type Service struct {
	Auth
	Profile
	Recipe
	RecipeOwnership
	RecipeSharing
	RecipePicture
	Encryption
	Category
	ShoppingList
}

type Dependencies struct {
	Repo                  *repository.Repository
	HashManager           hash.HashManager
	TokenManager          auth.TokenManager
	MailSender            mail.Sender
	MailConfig            config.MailConfig
	AccessTokenTTL        time.Duration
	RefreshTokenTTL       time.Duration
	CacheTTL              int64
	Environment           string
	Domain                string
	FirebaseImportEnabled bool
}

func NewService(dependencies Dependencies) *Service {

	mailService := service.NewMailService(dependencies.MailSender, dependencies.MailConfig)
	var firebaseService *service.FirebaseService = nil
	if dependencies.FirebaseImportEnabled {
		firebaseService = service.NewFirebaseService(dependencies.Repo.Migration, dependencies.Repo.Auth, dependencies.Repo.Profile,
			dependencies.Repo.Recipe, dependencies.Repo.RecipeOwnership, dependencies.Repo.Category, dependencies.Repo.ShoppingList)
	}

	return &Service{
		Auth: service.NewAuthService(dependencies.Repo.Auth, firebaseService, dependencies.HashManager, dependencies.TokenManager,
			dependencies.AccessTokenTTL, dependencies.RefreshTokenTTL, *mailService, dependencies.Domain),
		Profile:         service.NewProfileService(dependencies.Repo.Auth, dependencies.Repo.Profile, dependencies.Repo.File, dependencies.HashManager),
		Recipe:          service.NewRecipeService(dependencies.Repo.Recipe, dependencies.Repo.Category),
		RecipeOwnership: service.NewRecipeOwnershipService(dependencies.Repo.Recipe, dependencies.Repo.RecipeOwnership),
		RecipeSharing:   service.NewRecipeSharingService(dependencies.Repo.Recipe, dependencies.Repo.RecipeSharing),
		RecipePicture:   service.NewRecipePicturesService(dependencies.Repo.Recipe, dependencies.Repo.File),
		Encryption:      service.NewEncryptionService(dependencies.Repo.Encryption, dependencies.Repo.RecipeSharing, dependencies.Repo.Recipe, dependencies.Repo.File),
		Category:        service.NewCategoriesService(dependencies.Repo.Category),
		ShoppingList:    service.NewShoppingListService(dependencies.Repo.ShoppingList),
	}
}

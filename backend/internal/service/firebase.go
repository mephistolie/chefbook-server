package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/service/interface/repository"
	"github.com/mephistolie/chefbook-server/pkg/logger"
	"time"
)

const (
	oldUserBroccoins = 500
)

type FirebaseService struct {
	migrationRepo        repository.FirebaseMigration
	usersRepo            repository.Auth
	profileRepo          repository.Profile
	recipeRepo           repository.Recipe
	recipesOwnershipRepo repository.RecipeOwnership
	categoriesRepo       repository.Category
	shoppingListRepo     repository.ShoppingList
}

func NewFirebaseService(migrationRepo repository.FirebaseMigration, usersRepo repository.Auth, profileRepo repository.Profile, recipeRepo repository.Recipe,
	recipeOwnershipRepo repository.RecipeOwnership, categoriesRepo repository.Category, shoppingListRepo repository.ShoppingList) *FirebaseService {
	return &FirebaseService{
		migrationRepo:        migrationRepo,
		usersRepo:            usersRepo,
		profileRepo:          profileRepo,
		recipeRepo:           recipeRepo,
		recipesOwnershipRepo: recipeOwnershipRepo,
		categoriesRepo:       categoriesRepo,
		shoppingListRepo:     shoppingListRepo,
	}
}

func (s *FirebaseService) SignIn(credentials entity.Credentials) (entity.FirebaseProfile, error) {
	return s.migrationRepo.GetProfile(credentials)
}

func (s *FirebaseService) MigrateFromFirebase(credentials entity.Credentials, firebaseUser entity.FirebaseProfile) error {
	profile, err := s.migrationRepo.GetProfileData(firebaseUser)
	if err != nil {
		return failure.UnableImportFirebaseProfile
	}

	activationLink := uuid.New()
	userId, err := s.usersRepo.CreateUser(credentials, activationLink)
	if err != nil {
		return err
	}
	if err = s.usersRepo.ActivateProfile(activationLink); err != nil {
		logger.Warn("migration: error during activating user")
	}

	s.importProfileInfo(userId, profile.Profile)
	s.importFirebaseShoppingList(userId, profile.ShoppingList)
	categoriesIds := s.importCategories(userId, profile.Categories)
	s.importRecipes(userId, profile.Recipes, categoriesIds)

	return nil
}

func (s *FirebaseService) importProfileInfo(userId int, profile entity.FirebaseProfileInfo) {
	if profile.Username != nil {
		if err := s.profileRepo.SetUsername(userId, profile.Username); err != nil {
			logger.Warn("migration: error during setting name ", *profile.Username)
		}
	}

	if profile.CreationTimestamp != nil {
		if err := s.profileRepo.SetProfileCreationDate(userId, *profile.CreationTimestamp); err != nil {
			logger.Warn("migration: error during set profile creation date")
		}
	}

	if profile.IsPremium {
		lifetimePremium := time.Date(5000, 1, 1, 00, 00, 00, 00, time.UTC)
		err := s.profileRepo.SetPremiumDate(userId, lifetimePremium)
		if err != nil {
			logger.Warn("migration: error during activating premium")
		}
	}

	if err := s.profileRepo.IncreaseBroccoins(userId, oldUserBroccoins); err != nil {
		logger.Warn("migration: error during adding broccoins to old user")
	}
}

func (s *FirebaseService) importCategories(userId int, categories []entity.CategoryInput) map[string]int {
	categoriesIds := make(map[string]int)
	for _, category := range categories {
		categoryId, err := s.categoriesRepo.CreateCategory(category, userId)
		if err == nil {
			categoriesIds[category.Name] = categoryId
		}
	}

	return categoriesIds
}

func (s *FirebaseService) importRecipes(userId int, recipes []entity.FirebaseRecipe, categoriesIds map[string]int) {
	for _, recipe := range recipes {
		recipeId, err := s.recipesOwnershipRepo.CreateRecipe(recipe.Recipe, userId)
		if err != nil {
			logger.Error("migration: error during create recipe ", recipe.Recipe.Name)
			continue
		}

		var recipeCategoriesId []int
		for _, category := range recipe.Categories {
			recipeCategoriesId = append(recipeCategoriesId, categoriesIds[category])
		}
		_ = s.recipeRepo.SetRecipeCategories(recipeId, recipeCategoriesId, userId)

		_ = s.recipeRepo.SetRecipeFavourite(recipeId, recipe.IsFavourite, userId)
	}
}

func (s *FirebaseService) importFirebaseShoppingList(userId int, shoppingList entity.ShoppingList) {
	if err := s.shoppingListRepo.SetShoppingList(shoppingList, userId); err != nil {
		logger.Warn("migration: error during setting shopping list ")
	}
}

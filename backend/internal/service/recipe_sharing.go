package service

import (
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/service/interface/repository"
)

type RecipeSharingService struct {
	recipesRepo        repository.Recipe
	recipesSharingRepo repository.RecipeSharing
}

func NewRecipeSharingService(recipesRepo repository.Recipe, recipesSharingRepo repository.RecipeSharing) *RecipeSharingService {
	return &RecipeSharingService{
		recipesRepo:        recipesRepo,
		recipesSharingRepo: recipesSharingRepo,
	}
}

func (s *RecipeSharingService) GetUsersList(recipeId, userId string) ([]entity.ProfileInfo, error) {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return []entity.ProfileInfo{}, err
	}

	if ownerId != userId {
		return []entity.ProfileInfo{}, failure.NotOwner
	}

	return s.recipesSharingRepo.GetRecipeUserList(recipeId)
}

func (s *RecipeSharingService) GetUserPublicKey(recipeId, userId, requesterId string) (string, error) {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return "", err
	}

	if ownerId != requesterId {
		return "", failure.NotOwner
	}

	return s.recipesSharingRepo.GetUserRecipeKey(recipeId, userId)
}

func (s *RecipeSharingService) SetUserPublicKey(recipeId, userId string, userKey *string) error {
	err := s.recipesSharingRepo.SetUserPublicKeyLink(recipeId, userId, userKey)
	if err != nil {
		return err
	}

	if userKey == nil {
		_ = s.recipesSharingRepo.SetOwnerPrivateKeyLinkForUser(recipeId, userId, nil)
	}

	return nil
}

func (s *RecipeSharingService) GetOwnerPrivateKeyForUser(recipeId, userId string) (string, error) {
	return s.recipesSharingRepo.GetUserRecipeKey(recipeId, userId)
}

func (s *RecipeSharingService) SetOwnerPrivateKeyForUser(recipeId, userId, requesterId string, ownerKey *string) error {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return err
	}
	if ownerId != requesterId {
		return failure.NotOwner
	}

	return s.recipesSharingRepo.SetOwnerPrivateKeyLinkForUser(recipeId, userId, ownerKey)
}

func (s *RecipeSharingService) DeleteUserAccess(recipeId, userId, requesterId string) error {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return err
	}
	if recipe.OwnerId != requesterId {
		return failure.NotOwner
	}
	return s.recipesRepo.RemoveRecipeFromRecipeBook(recipeId, userId)
}

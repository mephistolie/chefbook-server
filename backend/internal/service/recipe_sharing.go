package service

import (
	"chefbook-server/internal/entity"
	"chefbook-server/internal/entity/failure"
	"chefbook-server/internal/service/interface/repository"
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

func (s *RecipeSharingService) GetUsersList(recipeId, userId int) ([]entity.ProfileInfo, error) {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return []entity.ProfileInfo{}, err
	}

	if ownerId != userId {
		return []entity.ProfileInfo{}, failure.NotOwner
	}

	return s.recipesSharingRepo.GetRecipeUserList(recipeId)
}

func (s *RecipeSharingService) GetUserPublicKey(recipeId, userId, requesterId int) (string, error) {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return "", err
	}

	if ownerId != requesterId {
		return "", failure.NotOwner
	}

	return s.recipesSharingRepo.GetUserRecipeKey(recipeId, userId)
}

func (s *RecipeSharingService) SetUserPublicKey(recipeId, userId int, userKey *string) error {
	err := s.recipesSharingRepo.SetUserPublicKeyLink(recipeId, userId, userKey)
	if err != nil {
		return err
	}

	if userKey == nil {
		_ = s.recipesSharingRepo.SetOwnerPrivateKeyLinkForUser(recipeId, userId, nil)
	}

	return nil
}

func (s *RecipeSharingService) GetOwnerPrivateKeyForUser(recipeId, userId int) (string, error) {
	return s.recipesSharingRepo.GetUserRecipeKey(recipeId, userId)
}

func (s *RecipeSharingService) SetOwnerPrivateKeyForUser(recipeId int, userId int, requesterId int, ownerKey *string) error {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return err
	}
	if ownerId != requesterId {
		return failure.NotOwner
	}

	return s.recipesSharingRepo.SetOwnerPrivateKeyLinkForUser(recipeId, userId, ownerKey)
}

func (s *RecipeSharingService) DeleteUserAccess(recipeId, userId, requesterId int) error {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return err
	}
	if recipe.OwnerId != requesterId {
		return failure.NotOwner
	}
	return s.recipesRepo.RemoveRecipeFromRecipeBook(recipeId, userId)
}

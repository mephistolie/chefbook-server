package service

import (
	"github.com/mephistolie/chefbook-server/internal/model"
	"github.com/mephistolie/chefbook-server/internal/repository"
)

type RecipeSharingService struct {
	recipesRepo           repository.RecipeCrud
	recipesSharingRepo    repository.RecipeSharing
}

func NewRecipeSharingService(recipesRepo repository.RecipeCrud, recipesSharingRepo repository.RecipeSharing) *RecipeSharingService {
	return &RecipeSharingService{
		recipesRepo:           recipesRepo,
		recipesSharingRepo:    recipesSharingRepo,
	}
}

func (s *RecipeSharingService) GetRecipeUserList(recipeId, userId int) ([]model.UserInfo, error) {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return []model.UserInfo{}, model.ErrRecipeNotFound
	}
	if recipe.OwnerId != userId {
		return []model.UserInfo{}, model.ErrAccessDenied
	}
	return s.recipesSharingRepo.GetRecipeUserList(recipeId)
}

func (s *RecipeSharingService) SetUserPublicKeyForRecipe(recipeId, userId int, userKey string) error {
	return s.recipesSharingRepo.SetUserPublicKeyForRecipe(recipeId, userId, userKey)
}

func (s *RecipeSharingService) SetUserPrivateKeyForRecipe(recipeId, userId int, userKey string) error {
	return s.recipesSharingRepo.SetUserPrivateKeyForRecipe(recipeId, userId, userKey)
}

func (s *RecipeSharingService) DeleteUserAccessToRecipe(recipeId, userId, requesterId int) error {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return model.ErrRecipeNotFound
	}
	if recipe.OwnerId != requesterId {
		return model.ErrAccessDenied
	}
	return s.recipesRepo.DeleteRecipeFromRecipeBook(recipeId, userId)
}

package service

import (
	"github.com/mephistolie/chefbook-server/internal/model"
	"github.com/mephistolie/chefbook-server/internal/repository"
)

type RecipeInteractionService struct {
	recipeInteractionRepo repository.RecipeInteraction
}

func NewRecipeInteractionService(recipeInteractionRepo repository.RecipeInteraction) *RecipeInteractionService {
	return &RecipeInteractionService{
		recipeInteractionRepo: recipeInteractionRepo,
	}
}

func (s *RecipeInteractionService) SetRecipeCategories(input model.RecipeCategoriesInput) error {
	return s.recipeInteractionRepo.SetRecipeCategories(input.Categories, input.RecipeId, input.UserId)
}

func (s *RecipeInteractionService) SetRecipeFavourite(input model.FavouriteRecipeInput) error {
	return s.recipeInteractionRepo.SetRecipeFavourite(input.RecipeId, input.UserId, input.Favourite)
}

func (s *RecipeInteractionService) SetRecipeLiked(input model.RecipeLikeInput) error {
	return s.recipeInteractionRepo.SetRecipeLiked(input.RecipeId, input.UserId, input.Liked)
}

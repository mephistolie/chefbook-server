package service

import (
	"chefbook-server/internal/entity"
	"chefbook-server/internal/entity/failure"
	"chefbook-server/internal/service/interface/repository"
)

type RecipeOwnershipService struct {
	recipeRepo    repository.Recipe
	ownershipRepo repository.RecipeOwnership
}

func NewRecipeOwnershipService(recipeRepo repository.Recipe, ownershipRepo repository.RecipeOwnership) *RecipeOwnershipService {
	return &RecipeOwnershipService{
		recipeRepo:    recipeRepo,
		ownershipRepo: ownershipRepo,
	}
}

func (s *RecipeOwnershipService) CreateRecipe(recipe entity.RecipeInput, userId int) (int, error) {
	return s.ownershipRepo.CreateRecipe(recipe, userId)
}

func (s *RecipeOwnershipService) UpdateRecipe(recipe entity.RecipeInput, recipeId, userId int) error {
	ownerId, err := s.recipeRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return err
	}
	if ownerId != userId {
		return failure.NotOwner
	}

	return s.ownershipRepo.UpdateRecipe(recipeId, recipe)
}

func (s *RecipeOwnershipService) DeleteRecipe(recipeId, userId int) error {
	ownerId, err := s.recipeRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return err
	}
	if ownerId != userId {
		return failure.NotOwner
	}

	return s.ownershipRepo.DeleteRecipe(recipeId)
}

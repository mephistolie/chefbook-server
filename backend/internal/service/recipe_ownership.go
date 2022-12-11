package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/service/interface/repository"
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

func (s *RecipeOwnershipService) CreateRecipe(recipe entity.RecipeInput, userId uuid.UUID) (uuid.UUID, error) {
	return s.ownershipRepo.CreateRecipe(recipe, userId)
}

func (s *RecipeOwnershipService) UpdateRecipe(recipe entity.RecipeInput, recipeId, userId uuid.UUID) error {
	ownerId, err := s.recipeRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return err
	}
	if ownerId != userId {
		return failure.NotOwner
	}

	return s.ownershipRepo.UpdateRecipe(recipeId, recipe)
}

func (s *RecipeOwnershipService) DeleteRecipe(recipeId, userId uuid.UUID) error {
	ownerId, err := s.recipeRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return err
	}
	if ownerId != userId {
		return failure.NotOwner
	}

	return s.ownershipRepo.DeleteRecipe(recipeId)
}

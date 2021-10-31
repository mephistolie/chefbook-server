package service

import (
	models2 "github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository"
)

type RecipesService struct {
	repo repository.Recipes
}

func NewRecipesService(repo repository.Recipes) *RecipesService {
	return &RecipesService{
		repo: repo,
	}
}

func (s *RecipesService) GetRecipesByUser(userId int) ([]models2.Recipe, error) {
	return s.repo.GetRecipesByUser(userId)
}

func (s *RecipesService) AddRecipe(recipe models2.Recipe) (int, error) {
	return s.repo.CreateRecipe(recipe)
}

func (s *RecipesService) GetRecipeById(recipeId, userId int) (models2.Recipe, error) {
	recipe, err := s.repo.GetRecipeById(recipeId, userId)
	if err != nil {
		return models2.Recipe{}, models2.ErrRecipeNotFound
	}
	return recipe, err
}

func (s *RecipesService) UpdateRecipe(recipe models2.Recipe, userId int) error {
	return s.repo.UpdateRecipe(recipe, userId)
}

func (s *RecipesService) DeleteRecipe(recipeId, userId int) error {
	ownerId, err := s.repo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return models2.ErrRecipeNotFound
	}
	if ownerId == userId {
		err = s.repo.DeleteRecipe(recipeId)
		return err
	} else {
		err = s.repo.DeleteRecipeLink(recipeId, userId)
		return err
	}
}

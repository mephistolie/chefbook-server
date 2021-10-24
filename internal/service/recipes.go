package service

import (
	"github.com/mephistolie/chefbook-server/internal/models"
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

func (s *RecipesService) AddRecipe(recipe models.Recipe) (int, error) {
	return s.repo.CreateRecipe(recipe)
}

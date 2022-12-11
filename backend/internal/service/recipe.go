package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/service/interface/repository"
	"strings"
)

type RecipeService struct {
	recipesRepo    repository.Recipe
	categoriesRepo repository.Category
}

func NewRecipeService(recipesRepo repository.Recipe, categoriesRepo repository.Category) *RecipeService {
	return &RecipeService{
		recipesRepo:    recipesRepo,
		categoriesRepo: categoriesRepo,
	}
}

func (s *RecipeService) GetRecipes(query entity.RecipesQuery, userId uuid.UUID) ([]entity.RecipeInfo, error) {
	recipes, err := s.recipesRepo.GetRecipes(query, userId)

	for i := range recipes {
		recipes[i].Categories = s.categoriesRepo.GetRecipeCategories(recipes[i].Id, userId)
		if recipes[i].OwnerId == userId {
			recipes[i].IsOwned = true
		}
	}
	return recipes, err
}

func (s *RecipeService) GetRecipe(recipeId, userId uuid.UUID) (entity.UserRecipe, error) {
	recipe, err := s.recipesRepo.GetRecipeWithUserFields(recipeId, userId)
	if err != nil {
		return entity.UserRecipe{}, err
	}

	if strings.ToLower(recipe.Visibility) == entity.VisibilityPrivate && recipe.OwnerId != userId {
		return entity.UserRecipe{}, failure.AccessDenied
	}

	recipe.Categories = s.categoriesRepo.GetRecipeCategories(recipeId, userId)
	if recipe.OwnerId == userId {
		recipe.IsOwned = true
	}

	return recipe, err
}

func (s *RecipeService) GetRandomRecipe(languages *[]string, userId uuid.UUID) (entity.UserRecipe, error) {
	return s.recipesRepo.GetRandomRecipe(languages, userId)
}

func (s *RecipeService) AddRecipeToRecipeBook(recipeId, userId uuid.UUID) error {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return err
	}

	if strings.ToLower(recipe.Visibility) == entity.VisibilityPrivate && recipe.OwnerId != userId {
		return failure.AccessDenied
	}

	err = s.recipesRepo.AddRecipeToRecipeBook(recipeId, userId)
	if err != nil {
		return failure.UnableAddRecipe
	}

	return nil
}

func (s *RecipeService) RemoveRecipeFromRecipeBook(recipeId, userId uuid.UUID) error {
	return s.recipesRepo.RemoveRecipeFromRecipeBook(recipeId, userId)
}

func (s *RecipeService) SetRecipeCategories(recipeId uuid.UUID, categories []uuid.UUID, userId uuid.UUID) error {
	return s.recipesRepo.SetRecipeCategories(recipeId, categories, userId)
}

func (s *RecipeService) SetRecipeFavourite(recipeId uuid.UUID, favourite bool, userId uuid.UUID) error {
	return s.recipesRepo.SetRecipeFavourite(recipeId, favourite, userId)
}

func (s *RecipeService) SetRecipeLikeStatus(recipeId uuid.UUID, favourite bool, userId uuid.UUID) error {
	return s.recipesRepo.SetRecipeLiked(recipeId, favourite, userId)
}

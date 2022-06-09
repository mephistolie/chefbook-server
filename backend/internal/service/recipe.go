package service

import (
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/service/interface/repository"
	"strings"
)

type RecipeService struct {
	recipesRepo            repository.Recipe
	categoriesRepo         repository.Category
}

func NewRecipeService(recipesRepo repository.Recipe, categoriesRepo repository.Category) *RecipeService {
	return &RecipeService{
		recipesRepo:            recipesRepo,
		categoriesRepo:         categoriesRepo,
	}
}

func (s *RecipeService)	GetRecipes(query entity.RecipesQuery, userId int) ([]entity.RecipeInfo, error) {
	recipes, err := s.recipesRepo.GetRecipes(query, userId)

	for i := range recipes {
		recipes[i].Categories= s.categoriesRepo.GetRecipeCategories(recipes[i].Id, userId)
		if recipes[i].OwnerId == userId {
			recipes[i].Owned = true
		}
	}
	return recipes, err
}

func (s *RecipeService) GetRecipe(recipeId, userId int) (entity.UserRecipe, error) {
	recipe, err := s.recipesRepo.GetRecipeWithUserFields(recipeId, userId)
	if err != nil {
		return entity.UserRecipe{}, err
	}

	if strings.ToLower(recipe.Visibility) == entity.VisibilityPrivate && recipe.OwnerId != userId {
		return entity.UserRecipe{}, failure.AccessDenied
	}

	recipe.Categories = s.categoriesRepo.GetRecipeCategories(recipeId, userId)
	if recipe.OwnerId == userId {
		recipe.Owned = true
	}

	return recipe, err
}

func (s *RecipeService) GetRandomRecipe(languages *[]string, userId int) (entity.UserRecipe, error) {
	return s.recipesRepo.GetRandomRecipe(languages, userId)
}

func (s *RecipeService) AddRecipeToRecipeBook(recipeId, userId int) error {
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

func (s *RecipeService) RemoveRecipeFromRecipeBook(recipeId, userId int) error {
	return s.recipesRepo.RemoveRecipeFromRecipeBook(recipeId, userId)
}

func (s *RecipeService) SetRecipeCategories(recipeId int, categories []int, userId int) error {
	return s.recipesRepo.SetRecipeCategories(recipeId, categories, userId)
}

func (s *RecipeService) SetRecipeFavourite(recipeId int, favourite bool, userId int) error {
	return s.recipesRepo.SetRecipeFavourite(recipeId, favourite, userId)
}

func (s *RecipeService) SetRecipeLikeStatus(recipeId int, favourite bool, userId int) error {
	return s.recipesRepo.SetRecipeLiked(recipeId, favourite, userId)
}

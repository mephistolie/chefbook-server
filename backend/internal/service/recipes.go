package service

import (
	"encoding/json"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"strings"
)

const (
	PRIVATE = "private"
	SHARED = "shared"
	PUBLIC = "public"
)

type RecipesService struct {
	repo repository.Recipes
}

func NewRecipesService(repo repository.Recipes) *RecipesService {
	return &RecipesService{
		repo: repo,
	}
}

func (s *RecipesService) GetRecipesByUser(userId int) ([]models.Recipe, error) {
	recipes, err := s.repo.GetRecipesByUser(userId)
	if recipes == nil {
		recipes = []models.Recipe{}
	}
	return recipes, err
}

func (s *RecipesService) AddRecipe(recipe models.Recipe) (int, error) {
	recipe, err := validateRecipe(recipe)
	if err != nil {
		return 0, models.ErrInvalidRecipeInput
	}
	return s.repo.CreateRecipe(recipe)
}

func (s *RecipesService) GetRecipeById(recipeId, userId int) (models.Recipe, error) {
	recipe, err := s.repo.GetRecipeById(recipeId, userId)
	if err != nil {
		return models.Recipe{}, models.ErrRecipeNotFound
	}
	return recipe, err
}

func (s *RecipesService) UpdateRecipe(recipe models.Recipe, userId int) error {
	recipe, err := validateRecipe(recipe)
	if err != nil {
		return models.ErrInvalidRecipeInput
	}
	return s.repo.UpdateRecipe(recipe, userId)
}

func (s *RecipesService) DeleteRecipe(recipeId, userId int) error {
	ownerId, err := s.repo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return models.ErrRecipeNotFound
	}
	if ownerId == userId {
		err = s.repo.DeleteRecipe(recipeId)
		return err
	} else {
		err = s.repo.DeleteRecipeLink(recipeId, userId)
		return err
	}
}

func (s *RecipesService) MarkRecipeFavourite(recipe models.FavouriteRecipeInput, userId int) error {
	return s.repo.MarkRecipeFavourite(recipe.RecipeId, userId, recipe.Favourite)
}

func validateRecipe(recipe models.Recipe) (models.Recipe, error) {
	var err error
	recipe.Ingredients, err = json.Marshal(recipe.Ingredients)
	if err != nil {
		return models.Recipe{}, err
	}

	recipe.Cooking, err = json.Marshal(recipe.Cooking)
	if err != nil {
		return models.Recipe{}, err
	}

	recipe.Visibility = strings.ToLower(recipe.Visibility)
	if recipe.Visibility != PUBLIC && recipe.Visibility != SHARED {
		recipe.Visibility = PRIVATE
	}

	if recipe.Servings < 1  {
		recipe.Servings = 1
	}

	if recipe.Time < 1  {
		recipe.Time = 15
	}

	return recipe, nil
}
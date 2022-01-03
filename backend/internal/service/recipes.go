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
	recipesRepo repository.Recipes
	categoriesRepo repository.Categories
}

func NewRecipesService(recipesRepo repository.Recipes, categoriesRepo repository.Categories) *RecipesService {
	return &RecipesService{
		recipesRepo: recipesRepo,
		categoriesRepo: categoriesRepo,
	}
}

func (s *RecipesService) GetRecipesByUser(userId int) ([]models.Recipe, error) {
	recipes, err := s.recipesRepo.GetRecipesByUser(userId)
	if recipes == nil {
		recipes = []models.Recipe{}
	}
	for i, _ := range recipes {
		recipes[i].Categories, _ = s.categoriesRepo.GetRecipeCategories(recipes[i].Id, userId)
	}
	return recipes, err
}

func (s *RecipesService) AddRecipe(recipe models.Recipe) (int, error) {
	recipe, err := validateRecipe(recipe)
	if err != nil {
		return 0, models.ErrInvalidRecipeInput
	}
	return s.recipesRepo.CreateRecipe(recipe)
}

func (s *RecipesService) GetRecipeById(recipeId, userId int) (models.Recipe, error) {
	recipe, err := s.recipesRepo.GetRecipeById(recipeId, userId)
	if err != nil {
		return models.Recipe{}, models.ErrRecipeNotFound
	}
	recipe.Categories, err = s.categoriesRepo.GetRecipeCategories(recipeId, userId)
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
	return s.recipesRepo.UpdateRecipe(recipe, userId)
}

func (s *RecipesService) DeleteRecipe(recipeId, userId int) error {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return models.ErrRecipeNotFound
	}
	if ownerId == userId {
		err = s.recipesRepo.DeleteRecipe(recipeId)
		return err
	} else {
		err = s.recipesRepo.DeleteRecipeLink(recipeId, userId)
		return err
	}
}

func (s *RecipesService) SetRecipeCategories(input models.RecipeCategoriesInput) error  {
	return s.recipesRepo.SetRecipeCategories(input.Categories, input.RecipeId, input.UserId)
}

func (s *RecipesService) MarkRecipeFavourite(input models.FavouriteRecipeInput) error {
	return s.recipesRepo.MarkRecipeFavourite(input.RecipeId, input.UserId, input.Favourite)
}

func (s *RecipesService) SetRecipeLike(input models.RecipeLikeInput) error  {
	return s.recipesRepo.SetRecipeLike(input.RecipeId, input.UserId, input.Liked)
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
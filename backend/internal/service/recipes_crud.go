package service

import (
	"encoding/json"
	"github.com/mephistolie/chefbook-server/internal/model"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"strings"
)

const (
	PRIVATE = "private"
	SHARED  = "shared"
	PUBLIC  = "public"
)

type RecipesService struct {
	recipesRepo            repository.RecipeCrud
	categoriesRepo         repository.Categories
}

func NewRecipesService(recipesRepo repository.RecipeCrud, categoriesRepo repository.Categories) *RecipesService {
	return &RecipesService{
		recipesRepo:            recipesRepo,
		categoriesRepo:         categoriesRepo,
	}
}

func (s *RecipesService) GetRecipesInfoByRequest(params model.RecipesRequestParams) ([]model.RecipeInfo, error) {
	recipes, err := s.recipesRepo.GetRecipesInfoByRequest(params)
	if recipes == nil {
		recipes = []model.RecipeInfo{}
	}
	for i, _ := range recipes {
		recipes[i].Categories, _ = s.categoriesRepo.GetRecipeCategories(recipes[i].Id, params.UserId)
		if recipes[i].OwnerId == params.UserId {
			recipes[i].Owned = true
		}
	}
	return recipes, err
}

func (s *RecipesService) CreateRecipe(recipe model.Recipe) (int, error) {
	recipe, err := validateRecipe(recipe)
	if err != nil {
		return 0, model.ErrInvalidInput
	}
	return s.recipesRepo.CreateRecipe(recipe)
}

func (s *RecipesService) AddRecipeToRecipeBook(recipeId, userId int) error {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return model.ErrRecipeNotFound
	}
	if strings.ToLower(recipe.Visibility) == "private" && recipe.OwnerId != userId {
		return model.ErrAccessDenied
	}
	err = s.recipesRepo.AddRecipeToRecipeBook(recipeId, userId)
	if err != nil {
		return model.ErrUnableAddRecipe
	}
	return nil
}

func (s *RecipesService) GetRecipeById(recipeId, userId int) (model.Recipe, error) {
	recipe, err := s.recipesRepo.GetRecipeWithUserFields(recipeId, userId)
	if err != nil {
		return model.Recipe{}, model.ErrRecipeNotFound
	}
	if strings.ToLower(recipe.Visibility) == "private" && recipe.OwnerId != userId {
		return model.Recipe{}, model.ErrAccessDenied
	}
	recipe.Categories, err = s.categoriesRepo.GetRecipeCategories(recipeId, userId)
	if err != nil {
		return model.Recipe{}, model.ErrRecipeNotFound
	}
	if recipe.OwnerId == userId {
		recipe.Owned = true
	}
	return recipe, err
}

func (s *RecipesService) GetRandomPublicRecipe(languages []string) (model.Recipe, error) {
	return s.recipesRepo.GetRandomPublicRecipe(languages)
}

func (s *RecipesService) UpdateRecipe(recipe model.Recipe) error {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipe.Id)
	if err != nil {
		return model.ErrRecipeNotFound
	}
	if ownerId != recipe.OwnerId {
		return model.ErrNotOwner
	}
	recipe, err = validateRecipe(recipe)
	if err != nil {
		return model.ErrInvalidInput
	}
	return s.recipesRepo.UpdateRecipe(recipe)
}

func (s *RecipesService) DeleteRecipe(recipeId, userId int) error {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return model.ErrRecipeNotFound
	}
	if ownerId == userId {
		err = s.recipesRepo.DeleteRecipe(recipeId)
		return err
	} else {
		err = s.recipesRepo.DeleteRecipeFromRecipeBook(recipeId, userId)
		return err
	}
}

func validateRecipe(recipe model.Recipe) (model.Recipe, error) {
	var err error
	recipe.Ingredients, err = json.Marshal(recipe.Ingredients)
	if err != nil {
		return model.Recipe{}, err
	}

	recipe.Cooking, err = json.Marshal(recipe.Cooking)
	if err != nil {
		return model.Recipe{}, err
	}

	recipe.Visibility = strings.ToLower(recipe.Visibility)
	if recipe.Visibility != PUBLIC && recipe.Visibility != SHARED {
		recipe.Visibility = PRIVATE
	}

	if len(recipe.Language) > 2 {
		recipe.Language = recipe.Language[0:2]
	} else if recipe.Language == "" {
		recipe.Language = "en"
	}

	if recipe.Servings < 1 {
		recipe.Servings = 1
	}

	if recipe.Time < 1 {
		recipe.Time = 15
	}

	return recipe, nil
}

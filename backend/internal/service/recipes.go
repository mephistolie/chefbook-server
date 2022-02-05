package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"github.com/mephistolie/chefbook-server/internal/repository/s3"
	"strings"
)

const (
	PRIVATE = "private"
	SHARED = "shared"
	PUBLIC = "public"
)

var imageTypes = map[string]interface{} {
	"image/jpeg": nil,
	"image/png": nil,
}

type RecipesService struct {
	recipesRepo repository.Recipes
	categoriesRepo repository.Categories
	filesRepo repository.Files
}

func NewRecipesService(recipesRepo repository.Recipes, categoriesRepo repository.Categories, filesRepo repository.Files) *RecipesService {
	return &RecipesService{
		recipesRepo: recipesRepo,
		categoriesRepo: categoriesRepo,
		filesRepo: filesRepo,
	}
}

func (s *RecipesService) GetRecipesInfoByRequest(params models.RecipesRequestParams) ([]models.RecipeInfo, error) {
	recipes, err := s.recipesRepo.GetRecipesInfoByRequest(params)
	if recipes == nil {
		recipes = []models.RecipeInfo{}
	}
	for i, _ := range recipes {
		recipes[i].Categories, _ = s.categoriesRepo.GetRecipeCategories(recipes[i].Id, params.UserId)
		if recipes[i].OwnerId == params.UserId {
			recipes[i].Owned = true
		}
	}
	return recipes, err
}

func (s *RecipesService) CreateRecipe(recipe models.Recipe) (int, error) {
	recipe, err := validateRecipe(recipe)
	if err != nil {
		return 0, models.ErrInvalidRecipeInput
	}
	return s.recipesRepo.CreateRecipe(recipe)
}

func (s *RecipesService) AddRecipeToRecipeBook(recipeId, userId int) error {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return models.ErrRecipeNotFound
	}
	if strings.ToLower(recipe.Visibility) == "private" && recipe.OwnerId != userId {
		return models.ErrAccessDenied
	}
	err = s.recipesRepo.AddRecipeLink(recipeId, userId)
	if err != nil {
		return models.ErrUnableAddRecipe
	}
	return nil
}

func (s *RecipesService) GetRecipeById(recipeId, userId int) (models.Recipe, error) {
	recipe, err := s.recipesRepo.GetRecipeWithUserFields(recipeId, userId)
	if err != nil {
		return models.Recipe{}, models.ErrRecipeNotFound
	}
	if strings.ToLower(recipe.Visibility) == "private" && recipe.OwnerId != userId {
		return models.Recipe{}, models.ErrAccessDenied
	}
	recipe.Categories, err = s.categoriesRepo.GetRecipeCategories(recipeId, userId)
	if err != nil {
		return models.Recipe{}, models.ErrRecipeNotFound
	}
	if recipe.OwnerId == userId {
		recipe.Owned = true
	}
	return recipe, err
}

func (s *RecipesService) UpdateRecipe(recipe models.Recipe, userId int) error {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipe.Id)
	if err != nil {
		return err
	}
	if ownerId != userId {
		return models.ErrNotOwner
	}
	recipe, err = validateRecipe(recipe)
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
		err = s.recipesRepo.DeleteRecipeFromRecipeBook(recipeId, userId)
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

func (s *RecipesService) UploadRecipePicture(ctx context.Context, recipeId, userId int, file *bytes.Reader, size int64, contentType string) (string, error) {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if !recipe.Encrypted {
		if _, ex := imageTypes[contentType]; !ex {
			return "", models.ErrFileTypeNotSupported
		}
	}
	if err != nil {
		return "", err
	}
	if recipe.OwnerId != userId {
		return "", models.ErrNotOwner
	}
	url, err := s.filesRepo.UploadRecipePicture(ctx, recipeId, s3.UploadInput{
		Name:        uuid.NewString(),
		File:        file,
		Size:        size,
		ContentType: contentType,
	})
	if err != nil {
		_ = s.filesRepo.DeleteFile(ctx, url)
		return "", err
	}
	return url, err
}

func (s *RecipesService) DeleteRecipePicture(ctx context.Context, recipeId, userId int, pictureName string) error  {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return err
	}
	if recipe.OwnerId != userId {
		return models.ErrNotOwner
	}
	url := s.filesRepo.GetRecipePictureLink(recipeId, pictureName)
	err = s.filesRepo.DeleteFile(ctx, url)
	return err
}

func (s *RecipesService) GetRecipeKey(recipeId, userId int) (string, error) {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return "", err
	}
	if recipe.OwnerId != userId {
		return "", models.ErrNotOwner
	}
	url, err := s.recipesRepo.GetRecipeKey(recipeId)
	if err != nil || url == "" {
		return "", models.ErrNoKey
	}
	return url, err
}

func (s *RecipesService) UploadRecipeKey(ctx context.Context, recipeId, userId int, file *bytes.Reader, size int64, contentType string) (string, error) {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return "", err
	}
	if recipe.OwnerId != userId {
		return "", models.ErrNotOwner
	}
	oldKey, err := s.recipesRepo.GetRecipeKey(recipeId)
	if err != nil {
		return "", err
	}
	url, err := s.filesRepo.UploadRecipeKey(ctx, recipeId, s3.UploadInput{
		Name:        uuid.NewString(),
		File:        file,
		Size:        size,
		ContentType: contentType,
	})
	err = s.recipesRepo.SetRecipeKey(recipeId, url)
	if err != nil {
		_ = s.filesRepo.DeleteFile(ctx, url)
		return "", err
	}
	if oldKey != "" {
		_ = s.filesRepo.DeleteFile(ctx, oldKey)
	}
	return url, err
}

func (s *RecipesService) DeleteRecipeKey(ctx context.Context, recipeId, userId int) error  {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return err
	}
	if recipe.OwnerId != userId {
		return models.ErrNotOwner
	}
	url, err := s.recipesRepo.GetRecipeKey(recipeId)
	if err != nil {
		return err
	}
	err = s.filesRepo.DeleteFile(ctx, url)
	if err != nil {
		return err
	}
	err = s.recipesRepo.SetRecipeKey(recipeId, "")
	return err
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

func (s *RecipesService) SetUserPublicKeyForRecipe(recipeId, userId int, userKey string) error {
	return s.recipesRepo.SetUserPublicKeyForRecipe(recipeId, userId, userKey)
}

func (s *RecipesService) SetUserPrivateKeyForRecipe(recipeId, userId int, userKey string) error {
	return s.recipesRepo.SetUserPrivateKeyForRecipe(recipeId, userId, userKey)
}

func (s *RecipesService) GetRecipeUserList(recipeId, userId int) ([]models.UserInfo, error) {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return []models.UserInfo{}, err
	}
	if recipe.OwnerId != userId {
		return []models.UserInfo{}, models.ErrAccessDenied
	}
	return s.recipesRepo.GetRecipeUserList(recipeId)
}

func (s *RecipesService) DeleteUserAccessToRecipe(recipeId, userId, requesterId int) error {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return err
	}
	if recipe.OwnerId != requesterId {
		return models.ErrAccessDenied
	}
	return s.recipesRepo.DeleteRecipeFromRecipeBook(recipeId, userId)
}
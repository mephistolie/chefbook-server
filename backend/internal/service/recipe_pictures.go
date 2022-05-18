package service

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/model"
	"github.com/mephistolie/chefbook-server/internal/repository"
)

var imageTypes = map[string]interface{}{
	"image/jpeg": nil,
	"image/png":  nil,
}

type RecipePicturesService struct {
	recipesRepo            repository.RecipeCrud
	filesRepo              repository.Files
}

func NewRecipePicturesService(recipesRepo repository.RecipeCrud, filesRepo repository.Files) *RecipePicturesService {
	return &RecipePicturesService{
		recipesRepo:            recipesRepo,
		filesRepo:              filesRepo,
	}
}

func (s *RecipePicturesService) GetRecipePictures(ctx context.Context, recipeId int, userId int) ([]string, error) {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return []string{}, model.ErrRecipeNotFound
	}
	if recipe.OwnerId != userId {
		return []string{}, model.ErrNotOwner
	}
	pictures := s.filesRepo.GetRecipePictures(ctx, recipeId)
	if pictures == nil {
		pictures = []string{}
	}
	return pictures, nil
}

func (s *RecipePicturesService) UploadRecipePicture(ctx context.Context, recipeId, userId int, file model.MultipartFileInfo) (string, error) {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return "", model.ErrRecipeNotFound
	}
	if !recipe.Encrypted {
		if _, ex := imageTypes[file.ContentType]; !ex {
			return "", model.ErrFileTypeNotSupported
		}
	}
	if recipe.OwnerId != userId {
		return "", model.ErrNotOwner
	}
	url, err := s.filesRepo.UploadRecipePicture(ctx, recipeId, file)
	if err != nil {
		_ = s.filesRepo.DeleteFile(ctx, url)
		return "", model.ErrUnableUploadFile
	}
	return url, err
}

func (s *RecipePicturesService) DeleteRecipePicture(ctx context.Context, recipeId, userId int, pictureName string) error {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return model.ErrRecipeNotFound
	}
	if recipe.OwnerId != userId {
		return model.ErrNotOwner
	}
	url := s.filesRepo.GetRecipePictureLink(recipeId, pictureName)
	err = s.filesRepo.DeleteFile(ctx, url)
	if err != nil {
		return model.ErrUnableDeleteRecipePicture
	}

	return nil
}

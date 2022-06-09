package service

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/service/interface/repository"
)

type RecipePicturesService struct {
	recipesRepo            repository.Recipe
	filesRepo              repository.File
}

func NewRecipePicturesService(recipesRepo repository.Recipe, filesRepo repository.File) *RecipePicturesService {
	return &RecipePicturesService{
		recipesRepo:            recipesRepo,
		filesRepo:              filesRepo,
	}
}

func (s *RecipePicturesService) GetRecipePictures(ctx context.Context, recipeId int, userId int) ([]string, error) {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return []string{}, err
	}
	if ownerId != userId {
		return []string{}, failure.NotOwner
	}

	pictures := s.filesRepo.GetRecipePictures(ctx, recipeId)
	if pictures == nil {
		pictures = []string{}
	}

	return pictures, nil
}

func (s *RecipePicturesService) UploadRecipePicture(ctx context.Context, recipeId, userId int, file entity.MultipartFile) (string, error) {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return "", err
	}
	if !recipe.IsEncrypted && !file.IsImage() {
		return "", failure.UnsupportedFileType
	}
	if recipe.OwnerId != userId {
		return "", failure.NotOwner
	}

	url, err := s.filesRepo.UploadRecipePicture(ctx, recipeId, file)
	if err != nil {
		_ = s.filesRepo.DeleteFile(ctx, url)
		return "", err
	}

	return url, nil
}

func (s *RecipePicturesService) DeleteRecipePicture(ctx context.Context, recipeId, userId int, pictureName string) error {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return err
	}
	if ownerId != userId {
		return failure.NotOwner
	}

	err = s.filesRepo.DeleteRecipePicture(ctx, recipeId, pictureName)
	if err != nil {
		return err
	}

	return nil
}

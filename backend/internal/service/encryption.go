package service

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/service/interface/repository"
)

type EncryptionService struct {
	encryptionRepo repository.Encryption
	sharingRepo    repository.RecipeSharing
	recipesRepo    repository.Recipe
	filesRepo      repository.File
}

func NewEncryptionService(encryptionRepo repository.Encryption, sharingRepo repository.RecipeSharing, recipesRepo repository.Recipe, filesRepo repository.File) *EncryptionService {
	return &EncryptionService{
		encryptionRepo: encryptionRepo,
		sharingRepo:    sharingRepo,
		recipesRepo:    recipesRepo,
		filesRepo:      filesRepo,
	}
}

func (s *EncryptionService) GetUserKeyLink(userId string) (string, error) {
	link, err := s.encryptionRepo.GetUserKeyLink(userId)
	if err != nil {
		return "", err
	}

	if link == nil {
		return "", failure.NoKey
	}

	return *link, nil
}

func (s *EncryptionService) UploadUserKey(ctx context.Context, userId string, file entity.MultipartFile) (string, error) {
	previousUrl, err := s.encryptionRepo.GetUserKeyLink(userId)
	if err != nil && err != failure.NoKey {
		return "", err
	}

	link, err := s.filesRepo.UploadUserKey(ctx, userId, file)
	if err != nil {
		return "", err
	}
	err = s.encryptionRepo.SetUserKeyLink(userId, &link)
	if err != nil {
		_ = s.filesRepo.DeleteFile(ctx, link)
		return "", err
	}

	if previousUrl != nil {
		_ = s.filesRepo.DeleteFile(ctx, *previousUrl)
	}

	return link, err
}

func (s *EncryptionService) DeleteUserKey(ctx context.Context, userId string) error {
	link, err := s.encryptionRepo.GetUserKeyLink(userId)
	if err != nil {
		return err
	}

	if link != nil {
		_ = s.filesRepo.DeleteFile(ctx, *link)
	}

	err = s.encryptionRepo.SetUserKeyLink(userId, nil)
	return err
}

func (s *EncryptionService) GetRecipeKey(recipeId, userId string) (string, error) {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return "", err
	}
	if ownerId != userId {
		return "", failure.NotOwner
	}

	link, err := s.encryptionRepo.GetRecipeKeyLink(recipeId)
	if err != nil {
		return "", err
	}

	if link == nil {
		return "", failure.NoKey
	}

	return *link, err
}

func (s *EncryptionService) UploadRecipeKey(ctx context.Context, recipeId, userId string, file entity.MultipartFile) (string, error) {
	ownerId, err := s.recipesRepo.GetRecipeOwnerId(recipeId)
	if err != nil {
		return "", err
	}
	if ownerId != userId {
		return "", failure.NotOwner
	}

	previousLink, err := s.encryptionRepo.GetRecipeKeyLink(recipeId)
	if err != nil && err != failure.NoKey {
		return "", err
	}

	url, err := s.filesRepo.UploadRecipeKey(ctx, recipeId, file)
	if err != nil {
		return "", err
	}
	err = s.encryptionRepo.SetRecipeKeyLink(recipeId, &url)
	if err != nil {
		_ = s.filesRepo.DeleteFile(ctx, url)
		return "", err
	}

	if previousLink != nil {
		_ = s.filesRepo.DeleteFile(ctx, *previousLink)
	}

	return url, err
}

func (s *EncryptionService) DeleteRecipeKey(ctx context.Context, recipeId, userId string) error {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return err
	}
	if recipe.OwnerId != userId {
		return failure.NotOwner
	}

	link, err := s.encryptionRepo.GetRecipeKeyLink(recipeId)
	if err != nil {
		return err
	}

	if link != nil {
		err = s.filesRepo.DeleteFile(ctx, *link)
		if err != nil {
			return err
		}
	}

	err = s.encryptionRepo.SetRecipeKeyLink(recipeId, nil)
	if err != nil {
		return err
	}

	return nil
}

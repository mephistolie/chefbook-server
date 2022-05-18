package service

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/model"
	"github.com/mephistolie/chefbook-server/internal/repository"
)

type EncryptionService struct {
	encryptionRepo repository.Encryption
	recipesRepo    repository.RecipeCrud
	filesRepo      repository.Files
}

func NewEncryptionService(encryptionRepo repository.Encryption, recipesRepo repository.RecipeCrud, filesRepo repository.Files) *EncryptionService {
	return &EncryptionService{
		encryptionRepo: encryptionRepo,
		recipesRepo:    recipesRepo,
		filesRepo:      filesRepo,
	}
}

func (s *EncryptionService) GetUserKeyLink(userId int) (string, error) {
	key, err := s.encryptionRepo.GetUserKey(userId)
	if err != nil || key == "" {
		return "", model.ErrNoKey
	}
	return key, err
}

func (s *EncryptionService) UploadUserKey(ctx context.Context, userId int, file model.MultipartFileInfo) (string, error) {
	key, err := s.encryptionRepo.GetUserKey(userId)
	if err != nil {
		return "", model.ErrNoKey
	}
	url, err := s.filesRepo.UploadUserKey(ctx, userId, file)
	if err != nil {
		return "", model.ErrUnableUploadFile
	}
	err = s.encryptionRepo.SetUserKey(userId, url)
	if err != nil {
		_ = s.filesRepo.DeleteFile(ctx, url)
		return "", model.ErrUnableSetUserKey
	}
	if key != "" {
		_ = s.filesRepo.DeleteFile(ctx, key)
	}
	return url, err
}

func (s *EncryptionService) DeleteUserKey(ctx context.Context, userId int) error {
	url, err := s.encryptionRepo.GetUserKey(userId)
	if err != nil {
		return model.ErrNoKey
	}
	err = s.filesRepo.DeleteFile(ctx, url)
	err = s.encryptionRepo.SetUserKey(userId, "")
	if err != nil {
		return model.ErrUnableSetUserKey
	}
	return err
}

func (s *EncryptionService) GetRecipeKey(recipeId, userId int) (string, error) {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return "", model.ErrRecipeNotFound
	}
	if recipe.OwnerId != userId {
		return "", model.ErrNotOwner
	}
	url, err := s.encryptionRepo.GetRecipeKey(recipeId)
	if err != nil || url == "" {
		return "", model.ErrNoKey
	}
	return url, err
}

func (s *EncryptionService) UploadRecipeKey(ctx context.Context, recipeId, userId int, file model.MultipartFileInfo) (string, error) {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return "", model.ErrRecipeNotFound
	}
	if recipe.OwnerId != userId {
		return "", model.ErrNotOwner
	}
	oldKey, err := s.encryptionRepo.GetRecipeKey(recipeId)
	if err != nil {
		return "", model.ErrNoKey
	}
	url, err := s.filesRepo.UploadRecipeKey(ctx, recipeId, file)
	if err != nil {
		return "", model.ErrUnableUploadFile
	}
	err = s.encryptionRepo.SetRecipeKey(recipeId, url)
	if err != nil {
		_ = s.filesRepo.DeleteFile(ctx, url)
		return "", model.ErrUnableDeleteRecipeKey
	}
	if oldKey != "" {
		_ = s.filesRepo.DeleteFile(ctx, oldKey)
	}
	return url, err
}

func (s *EncryptionService) DeleteRecipeKey(ctx context.Context, recipeId, userId int) error {
	recipe, err := s.recipesRepo.GetRecipe(recipeId)
	if err != nil {
		return model.ErrRecipeNotFound
	}
	if recipe.OwnerId != userId {
		return model.ErrNotOwner
	}
	url, err := s.encryptionRepo.GetRecipeKey(recipeId)
	if err != nil {
		return model.ErrNoKey
	}
	err = s.filesRepo.DeleteFile(ctx, url)
	if err != nil {
		return model.ErrUnableDeleteRecipeKey
	}
	err = s.encryptionRepo.SetRecipeKey(recipeId, "")
	if err != nil {
		return model.ErrUnableDeleteRecipeKey
	}

	return nil
}

package repository

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type File interface {
	UploadAvatar(ctx context.Context, userId string, input entity.MultipartFile) (string, error)
	UploadUserKey(ctx context.Context, userId string, input entity.MultipartFile) (string, error)
	GetRecipePictures(ctx context.Context, recipeId string) []string
	UploadRecipePicture(ctx context.Context, recipeId string, input entity.MultipartFile) (string, error)
	DeleteRecipePicture(ctx context.Context, recipeId string, pictureName string) error
	UploadRecipeKey(ctx context.Context, recipeId string, input entity.MultipartFile) (string, error)
	DeleteFile(ctx context.Context, url string) error
}

package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type File interface {
	UploadAvatar(ctx context.Context, userId uuid.UUID, input entity.MultipartFile) (string, error)
	UploadUserKey(ctx context.Context, userId uuid.UUID, input entity.MultipartFile) (string, error)
	GetRecipePictures(ctx context.Context, recipeId uuid.UUID) []string
	UploadRecipePicture(ctx context.Context, recipeId uuid.UUID, input entity.MultipartFile) (string, error)
	DeleteRecipePicture(ctx context.Context, recipeId uuid.UUID, pictureName string) error
	UploadRecipeKey(ctx context.Context, recipeId uuid.UUID, input entity.MultipartFile) (string, error)
	DeleteFile(ctx context.Context, url string) error
}

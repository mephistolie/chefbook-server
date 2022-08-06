package repository

import (
	"chefbook-server/internal/entity"
	"context"
)

type File interface {
	UploadAvatar(ctx context.Context, userId int, input entity.MultipartFile) (string, error)
	UploadUserKey(ctx context.Context, userId int, input entity.MultipartFile) (string, error)
	GetRecipePictures(ctx context.Context, recipeId int) []string
	UploadRecipePicture(ctx context.Context, recipeId int, input entity.MultipartFile) (string, error)
	DeleteRecipePicture(ctx context.Context, recipeId int, pictureName string) error
	UploadRecipeKey(ctx context.Context, recipeId int, input entity.MultipartFile) (string, error)
	DeleteFile(ctx context.Context, url string) error
}

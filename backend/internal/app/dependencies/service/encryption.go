package service

import (
	"chefbook-server/internal/entity"
	"context"
)

type Encryption interface {
	GetUserKeyLink(userId int) (string, error)
	UploadUserKey(ctx context.Context, userId int, file entity.MultipartFile) (string, error)
	DeleteUserKey(ctx context.Context, userId int) error
	GetRecipeKey(recipeId, userId int) (string, error)
	UploadRecipeKey(ctx context.Context, recipeId, userId int, file entity.MultipartFile) (string, error)
	DeleteRecipeKey(ctx context.Context, recipeId, userId int) error
}

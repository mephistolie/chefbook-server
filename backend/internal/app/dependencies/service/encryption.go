package service

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Encryption interface {
	GetUserKeyLink(userId int) (string, error)
	UploadUserKey(ctx context.Context, userId int, file entity.MultipartFile) (string, error)
	DeleteUserKey(ctx context.Context, userId int) error
	GetRecipeKey(recipeId, userId int) (string, error)
	UploadRecipeKey(ctx context.Context, recipeId, userId int, file entity.MultipartFile) (string, error)
	DeleteRecipeKey(ctx context.Context, recipeId, userId int) error
}

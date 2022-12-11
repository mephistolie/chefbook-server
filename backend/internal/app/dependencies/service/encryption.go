package service

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Encryption interface {
	GetUserKeyLink(userId string) (string, error)
	UploadUserKey(ctx context.Context, userId string, file entity.MultipartFile) (string, error)
	DeleteUserKey(ctx context.Context, userId string) error
	GetRecipeKey(recipeId, userId string) (string, error)
	UploadRecipeKey(ctx context.Context, recipeId, userId string, file entity.MultipartFile) (string, error)
	DeleteRecipeKey(ctx context.Context, recipeId, userId string) error
}

package service

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Encryption interface {
	GetUserKeyLink(userId uuid.UUID) (string, error)
	UploadUserKey(ctx context.Context, userId uuid.UUID, file entity.MultipartFile) (string, error)
	DeleteUserKey(ctx context.Context, userId uuid.UUID) error
	GetRecipeKey(recipeId, userId uuid.UUID) (string, error)
	UploadRecipeKey(ctx context.Context, recipeId, userId uuid.UUID, file entity.MultipartFile) (string, error)
	DeleteRecipeKey(ctx context.Context, recipeId, userId uuid.UUID) error
}

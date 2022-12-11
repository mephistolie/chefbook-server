package service

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Profile interface {
	GetProfile(userId string) (entity.Profile, error)
	ChangePassword(userId string, oldPassword string, newPassword string) error
	SetUsername(userId string, username *string) error
	UploadAvatar(ctx context.Context, userId string, file entity.MultipartFile) (string, error)
	DeleteAvatar(ctx context.Context, userId string) error
}

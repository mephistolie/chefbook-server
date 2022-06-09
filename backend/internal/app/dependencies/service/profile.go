package service

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Profile interface {
	GetProfile(userId int) (entity.Profile, error)
	ChangePassword(userId int, oldPassword string, newPassword string) error
	SetUsername(userId int, username *string) error
	UploadAvatar(ctx context.Context, userId int, file entity.MultipartFile) (string, error)
	DeleteAvatar(ctx context.Context, userId int) error
}

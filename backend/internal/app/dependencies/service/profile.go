package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Profile interface {
	GetProfile(userId uuid.UUID) (entity.Profile, error)
	ChangePassword(userId uuid.UUID, oldPassword string, newPassword string) error
	SetUsername(userId uuid.UUID, username *string) error
	UploadAvatar(ctx context.Context, userId uuid.UUID, file entity.MultipartFile) (string, error)
	DeleteAvatar(ctx context.Context, userId uuid.UUID) error
}

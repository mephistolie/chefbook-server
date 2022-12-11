package service

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/service/interface/repository"
	"github.com/mephistolie/chefbook-server/pkg/hash"
)

type ProfileService struct {
	authRepo    repository.Auth
	profileRepo repository.Profile
	filesRepo   repository.File

	hashManager hash.HashManager
}

func NewProfileService(usersRepo repository.Auth, profileRepo repository.Profile, filesRepo repository.File, hashManager hash.HashManager) *ProfileService {
	return &ProfileService{
		authRepo:    usersRepo,
		profileRepo: profileRepo,
		filesRepo:   filesRepo,
		hashManager: hashManager,
	}
}

func (s *ProfileService) GetProfile(userId string) (entity.Profile, error) {
	return s.authRepo.GetUserById(userId)
}

func (s *ProfileService) ChangePassword(userId string, oldPassword string, newPassword string) error {
	profile, err := s.authRepo.GetUserById(userId)
	if err != nil {
		return err
	}

	if err = s.hashManager.ValidateByHash(oldPassword, profile.Password); err != nil {
		return failure.InvalidCredentials
	}

	newHashedPassword, err := s.hashManager.Hash(newPassword)
	if err != nil {
		return failure.Unknown
	}

	return s.authRepo.ChangePassword(userId, newHashedPassword)
}

func (s *ProfileService) SetUsername(userId string, username *string) error {
	return s.profileRepo.SetUsername(userId, username)
}

func (s *ProfileService) UploadAvatar(ctx context.Context, userId string, file entity.MultipartFile) (string, error) {
	user, err := s.authRepo.GetUserById(userId)
	if err != nil {
		return "", err
	}

	url, err := s.filesRepo.UploadAvatar(ctx, userId, file)
	if err != nil {
		return "", failure.UnableUploadFile
	}
	err = s.profileRepo.SetAvatarLink(userId, &url)
	if err != nil {
		_ = s.filesRepo.DeleteFile(ctx, url)
		return "", failure.UnableSetAvatar
	}

	if user.Avatar != nil {
		_ = s.filesRepo.DeleteFile(ctx, *user.Avatar)
	}

	return url, nil
}

func (s *ProfileService) DeleteAvatar(ctx context.Context, userId string) error {
	user, err := s.authRepo.GetUserById(userId)
	if err != nil {
		return err
	}

	err = s.filesRepo.DeleteFile(ctx, *user.Avatar)
	if err != nil {
		return err
	}
	err = s.profileRepo.SetAvatarLink(userId, nil)
	if err != nil {
		return err
	}

	return nil
}

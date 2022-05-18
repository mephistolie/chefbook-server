package service

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/model"
	"github.com/mephistolie/chefbook-server/internal/repository"
)

type ProfileService struct {
	usersRepo   repository.Auth
	profileRepo repository.Profile
	filesRepo   repository.Files
}

func NewUsersService(usersRepo repository.Auth, profileRepo repository.Profile, filesRepo repository.Files) *ProfileService {
	return &ProfileService{
		usersRepo:   usersRepo,
		profileRepo: profileRepo,
		filesRepo:   filesRepo,
	}
}

func (s *ProfileService) GetUserInfo(userId int) (model.User, error) {
	user, err := s.usersRepo.GetUserById(userId)
	if err != nil {
		return model.User{}, model.ErrUserNotFound
	}
	return user, nil
}

func (s *ProfileService) SetUsername(userId int, username string) error {
	return s.profileRepo.SetUsername(userId, username)
}

func (s *ProfileService) UploadAvatar(ctx context.Context, userId int, file model.MultipartFileInfo) (string, error) {
	user, err := s.usersRepo.GetUserById(userId)
	if err != nil {
		return "", model.ErrUserIdNotFound
	}
	url, err := s.filesRepo.UploadAvatar(ctx, userId, file)
	if err != nil {
		return "", model.ErrUnableUploadFile
	}
	err = s.profileRepo.SetAvatar(userId, url)
	if err != nil {
		_ = s.filesRepo.DeleteFile(ctx, url)
		return "", model.ErrUnableSetAvatar
	}
	if user.Avatar.String != "" {
		_ = s.filesRepo.DeleteFile(ctx, user.Avatar.String)
	}
	return url, nil
}

func (s *ProfileService) DeleteAvatar(ctx context.Context, userId int) error {
	user, err := s.usersRepo.GetUserById(userId)
	if err != nil {
		return model.ErrUserIdNotFound
	}
	err = s.filesRepo.DeleteFile(ctx, user.Avatar.String)
	if err != nil {
		return model.ErrUnableDeleteAvatar
	}
	err = s.profileRepo.SetAvatar(userId, "")
	if err != nil {
		return model.ErrUnableDeleteAvatar
	}
	return nil
}

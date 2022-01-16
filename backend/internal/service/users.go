package service

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"github.com/mephistolie/chefbook-server/internal/repository/s3"
)

type UsersService struct {
	usersRepo repository.Users
	filesRepo repository.Files
}

func NewUsersService(usersRepo repository.Users, filesRepo repository.Files) *UsersService {
	return &UsersService{
		usersRepo: usersRepo,
		filesRepo: filesRepo,
	}
}

func (s *UsersService) GetUserInfo(userId int) (models.User, error) {
	user, err := s.usersRepo.GetUserById(userId)
	if err != nil {
		return models.User{}, models.ErrUserNotFound
	}
	return user, nil
}

func (s *UsersService) SetUserName(userId int, username string) error  {
	return s.usersRepo.SetUserName(userId, username)
}

func (s *UsersService) UploadAvatar(ctx context.Context, userId int, file *bytes.Reader, size int64, contentType string) (string, error) {
	user, err := s.usersRepo.GetUserById(userId)
	if err != nil {
		return "", err
	}
	url, err := s.filesRepo.UploadAvatar(ctx, userId, s3.UploadInput{
		Name:        uuid.NewString(),
		File:        file,
		Size:        size,
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}
	err = s.usersRepo.SetUserAvatar(userId, url)
	if err != nil {
		_ = s.filesRepo.DeleteFile(ctx, url)
		return "", err
	}
	if user.Avatar.String != "" {
		_ = s.filesRepo.DeleteFile(ctx, user.Avatar.String)
	}
	return url, nil
}

func (s *UsersService) DeleteAvatar(ctx context.Context, userId int) error  {
	user, err := s.usersRepo.GetUserById(userId)
	if err != nil {
		return err
	}
	err = s.filesRepo.DeleteFile(ctx, user.Avatar.String)
	if err != nil {
		return err
	}
	err = s.usersRepo.SetUserAvatar(userId, "")
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersService) GetUserKey(userId int) (string, error) {
	key, err := s.usersRepo.GetUserKey(userId)
	if err != nil || key == "" {
		return "", models.ErrNoKey
	}
	return key, err
}

func (s *UsersService) UploadUserKey(ctx context.Context, userId int, file *bytes.Reader, size int64, contentType string) (string, error) {
	key, err := s.usersRepo.GetUserKey(userId)
	if err != nil {
		return "", err
	}
	url, err := s.filesRepo.UploadUserKey(ctx, userId, s3.UploadInput{
		Name:        uuid.NewString(),
		File:        file,
		Size:        size,
		ContentType: contentType,
	})
	err = s.usersRepo.SetUserKey(userId, url)
	if err != nil {
		_ = s.filesRepo.DeleteFile(ctx, url)
		return "", err
	}
	if key != "" {
		_ = s.filesRepo.DeleteFile(ctx, key)
	}
	return url, err
}

func (s *UsersService) DeleteUserKey(ctx context.Context, userId int) error  {
	url, err := s.usersRepo.GetUserKey(userId)
	if err != nil {
		return err
	}
	err = s.filesRepo.DeleteFile(ctx, url)
	return err
}
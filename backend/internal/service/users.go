package service

import (
	"bytes"
	"context"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"github.com/mephistolie/chefbook-server/internal/repository/s3"
	"strconv"
)

type UsersService struct {
	repo repository.Repository
}

func NewUsersService(repo repository.Repository) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (s *UsersService) GetUserInfo(userId int) (models.User, error) {
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return models.User{}, models.ErrUserNotFound
	}
	return user, nil
}

func (s *UsersService) SetUserName(userId int, username string) error  {
	return s.SetUserName(userId, username)
}

func (s *UsersService) UploadAvatar(ctx context.Context, userId int, file *bytes.Reader, size int64, contentType string) (string, error) {
	url, err := s.repo.Files.UploadAvatar(ctx, s3.UploadInput{
		Name:        strconv.Itoa(userId),
		File:        file,
		Size:        size,
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}
	err = s.repo.Users.SetUserAvatar(userId, url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *UsersService) DeleteAvatar(ctx context.Context, userId int) error  {
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return err
	}
	err = s.repo.DeleteAvatar(ctx, user.Avatar.String)
	if err != nil {
		return err
	}
	err = s.repo.Users.SetUserAvatar(userId, "")
	if err != nil {
		return err
	}
	return nil
}

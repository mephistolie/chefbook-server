package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mephistolie/chefbook-server/pkg/logger"
	"github.com/minio/minio-go/v7"
	"strings"
)

const chefBookBucket = "chefbook-storage"
const usersDir = "users"
const avatarsDir = "avatars"
const keysDir = "keys"
const recipesDir = "recipes"
const imagesDir = "images"

type UploadInput struct {
	File	*bytes.Reader
	Name	string
	Size	int64
	ContentType string
}

type AWSFileManager struct {
	client *minio.Client
}

func NewAWSFileManager(client *minio.Client) *AWSFileManager {
	return &AWSFileManager{
		client: client,
	}
}

func (r *AWSFileManager) UploadAvatar(ctx context.Context, userId int, input UploadInput) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType: input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	filePath := fmt.Sprintf("%s/%d/%s/%s.jpg", usersDir, userId, avatarsDir, input.Name)
	_, err := r.client.PutObject(ctx, chefBookBucket, filePath, input.File, input.Size, opts)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath), nil
}

func (r *AWSFileManager) UploadUserKey(ctx context.Context, userId int, input UploadInput) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType: input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	filePath := fmt.Sprintf("%s/%d/%s/%s", usersDir, userId, keysDir, input.Name)
	_, err := r.client.PutObject(ctx, chefBookBucket, filePath, input.File, input.Size, opts)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath), nil
}

func (r *AWSFileManager) UploadRecipePicture(ctx context.Context, recipeId int, input UploadInput) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType: input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	filePath := fmt.Sprintf("%s/%d/%s/%s.jpg", recipesDir, recipeId, imagesDir, input.Name)
	_, err := r.client.PutObject(ctx, chefBookBucket, filePath, input.File, input.Size, opts)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath), nil
}

func (r *AWSFileManager) UploadRecipeKey(ctx context.Context, recipeId int, input UploadInput) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType: input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	filePath := fmt.Sprintf("%s/%d/%s/%s", recipesDir, recipeId, keysDir, input.Name)
	_, err := r.client.PutObject(ctx, chefBookBucket, filePath, input.File, input.Size, opts)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath), nil
}

func (r *AWSFileManager) DeleteFile(ctx context.Context, url string) error {
	opts := minio.RemoveObjectOptions{ ForceDelete: true }
	filePath := strings.ReplaceAll(url, fmt.Sprintf("%s/%s/", r.client.EndpointURL().String(), chefBookBucket), "")
	err := r.client.RemoveObject(ctx, chefBookBucket, filePath, opts)
	logger.Error(err)
	return err
}

func (r *AWSFileManager) GetRecipePictureLink(recipeId int, pictureName string) string {
	filePath := fmt.Sprintf("%s/%d/%s/%s", recipesDir, recipeId, imagesDir, pictureName)
	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath)
}

func (r *AWSFileManager) GetRecipeKeysLink(recipeId int, pictureName string) string {
	filePath := fmt.Sprintf("%s/%d/%s/%s", recipesDir, recipeId, keysDir, pictureName)
	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath)
}
package s3

import (
	"context"
	"fmt"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/minio/minio-go/v7"
	"strings"
)

const (
	chefBookBucket = "chefbook-storage"
	usersDir       = "users"
	avatarsDir     = "avatars"
	keysDir        = "keys"
	recipesDir     = "recipes"
	imagesDir      = "images"

	xAmzAcl    = "x-amz-acl"
	publicRead = "public-read"
)

type AWSFileManager struct {
	client *minio.Client
}

func NewAWSFileManager(client *minio.Client) *AWSFileManager {
	return &AWSFileManager{
		client: client,
	}
}

func (r *AWSFileManager) UploadAvatar(ctx context.Context, userId uuid.UUID, input entity.MultipartFile) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{xAmzAcl: publicRead},
	}

	filePath := fmt.Sprintf("%s/%d/%s/%s", usersDir, userId, avatarsDir, input.Name)
	_, err := r.client.PutObject(ctx, chefBookBucket, filePath, input.Content, input.Size, opts)
	if err != nil {
		return "", failure.UnableUploadFile
	}

	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath), nil
}

func (r *AWSFileManager) UploadUserKey(ctx context.Context, userId uuid.UUID, input entity.MultipartFile) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{xAmzAcl: publicRead},
	}

	filePath := fmt.Sprintf("%s/%d/%s/%s", usersDir, userId, keysDir, input.Name)
	_, err := r.client.PutObject(ctx, chefBookBucket, filePath, input.Content, input.Size, opts)
	if err != nil {
		return "", failure.UnableUploadFile
	}

	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath), nil
}

func (r *AWSFileManager) GetRecipePictures(ctx context.Context, recipeId uuid.UUID) []string {
	picturesPath := fmt.Sprintf("%s/%d/%s", recipesDir, recipeId, imagesDir)
	var objects []string
	for object := range r.client.ListObjects(ctx, chefBookBucket, minio.ListObjectsOptions{Prefix: picturesPath, Recursive: true}) {
		objects = append(objects, fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, object.Key))
	}
	return objects
}

func (r *AWSFileManager) UploadRecipePicture(ctx context.Context, recipeId uuid.UUID, input entity.MultipartFile) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{xAmzAcl: publicRead},
	}

	filePath := fmt.Sprintf("%s/%d/%s/%s", recipesDir, recipeId, imagesDir, input.Name)
	_, err := r.client.PutObject(ctx, chefBookBucket, filePath, input.Content, input.Size, opts)
	if err != nil {
		return "", failure.UnableUploadFile
	}

	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath), nil
}

func (r *AWSFileManager) DeleteRecipePicture(ctx context.Context, recipeId uuid.UUID, pictureName string) error {
	return r.DeleteFile(ctx, r.getRecipePictureLink(recipeId, pictureName))
}

func (r *AWSFileManager) UploadRecipeKey(ctx context.Context, recipeId uuid.UUID, input entity.MultipartFile) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{xAmzAcl: publicRead},
	}

	filePath := fmt.Sprintf("%s/%d/%s/%s", recipesDir, recipeId, keysDir, input.Name)
	_, err := r.client.PutObject(ctx, chefBookBucket, filePath, input.Content, input.Size, opts)
	if err != nil {
		return "", failure.UnableUploadFile
	}

	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath), nil
}

func (r *AWSFileManager) DeleteFile(ctx context.Context, url string) error {
	opts := minio.RemoveObjectOptions{ForceDelete: true}
	filePath := strings.ReplaceAll(url, fmt.Sprintf("%s/%s/", r.client.EndpointURL().String(), chefBookBucket), "")
	if err := r.client.RemoveObject(ctx, chefBookBucket, filePath, opts); err != nil {
		return failure.UnableDeleteFile
	}
	return nil
}

func (r *AWSFileManager) getRecipePictureLink(recipeId uuid.UUID, pictureName string) string {
	filePath := fmt.Sprintf("%s/%d/%s/%s", recipesDir, recipeId, imagesDir, pictureName)
	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath)
}

func (r *AWSFileManager) getRecipeKeysLink(recipeId uuid.UUID, pictureName string) string {
	filePath := fmt.Sprintf("%s/%d/%s/%s", recipesDir, recipeId, keysDir, pictureName)
	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath)
}

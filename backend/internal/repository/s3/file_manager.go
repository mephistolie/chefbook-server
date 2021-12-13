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
const imagesDirectory = "images"
const avatarsDirectory = imagesDirectory + "/avatars"

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

func (r *AWSFileManager) UploadAvatar(ctx context.Context, input UploadInput) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType: input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	filePath := fmt.Sprintf("%s/%s.jpg", avatarsDirectory, input.Name)
	_, err := r.client.PutObject(ctx, chefBookBucket, filePath, input.File, input.Size, opts)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", r.client.EndpointURL(), chefBookBucket, filePath), nil
}

func (r *AWSFileManager) DeleteAvatar(ctx context.Context, url string) error {
	opts := minio.RemoveObjectOptions{ ForceDelete: true }
	filePath := strings.ReplaceAll(url, fmt.Sprintf("%s/%s/", r.client.EndpointURL().String(), chefBookBucket), "")
	logger.Error(filePath)
	err := r.client.RemoveObject(ctx, chefBookBucket, filePath, opts)
	logger.Error(err)
	return err
}
package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"mime/multipart"
	"net/http"
)



type FileMiddleware struct {
}

func NewFile() *FileMiddleware {
	return &FileMiddleware{}
}

func (h *FileMiddleware) GetFileWithMaxSize(c *gin.Context, maxSize int64) (entity.MultipartFile, error) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return entity.MultipartFile{}, failure.InvalidFileSize
	}

	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	buffer := make([]byte, header.Size)
	_, err = file.Read(buffer)
	if err != nil {
		return entity.MultipartFile{}, failure.InvalidFileSize
	}

	return entity.MultipartFile{
		Name:        uuid.NewString(),
		Content:     bytes.NewReader(buffer),
		Size:        header.Size,
		ContentType: http.DetectContentType(buffer),
	}, nil
}

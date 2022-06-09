package entity

import (
	"bytes"
)

var imageTypes = map[string]interface{}{
	"image/jpeg": nil,
	"image/png":  nil,
}

type MultipartFile struct {
	Name        string
	Content     *bytes.Reader
	Size        int64
	ContentType string
}

func (h *MultipartFile) IsImage() bool {
	_, ok := imageTypes[h.ContentType]
	return ok
}
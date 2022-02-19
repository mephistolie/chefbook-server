package model

import "bytes"

type MultipartFileInfo struct {
	File	*bytes.Reader
	Name	string
	Size	int64
	ContentType string
}

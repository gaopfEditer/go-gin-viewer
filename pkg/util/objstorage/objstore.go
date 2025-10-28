package objstorage

import (
	"io"
)

// Client 定义通用接口
type Client interface {
	UploadObject(objectName string, reader io.Reader) error
	DownloadObject(objectName string) (io.Reader, error)
	DeleteObject(objectName string) error
	ListObjects(prefix string) ([]string, error)
	GetTemporaryURL(objectName string) (string, error)
}

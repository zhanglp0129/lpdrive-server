package repository

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/zhanglp0129/lpdrive-server/common/constant/minioconstant"
	"github.com/zhanglp0129/lpdrive-server/config"
	"io"
)

// PutObject 上传文件
func PutObject(sha256 string, data io.Reader, size int64) error {
	_, err := MC.PutObject(context.Background(),
		config.C.Minio.BucketName,
		sha256, data, size,
		"", sha256, minio.PutObjectOptions{})
	return err
}

// ReadObject 读取文件
func ReadObject(sha256 string) (io.ReadCloser, error) {
	reader, _, _, err := MC.GetObject(context.Background(),
		config.C.Minio.BucketName,
		sha256, minio.StatObjectOptions{})
	return reader, err
}

// GetTempUploadObjectName 获取临时上传对象名
func GetTempUploadObjectName(uuid string) string {
	return minioconstant.UploadTempPrefix + uuid
}

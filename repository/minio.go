package repository

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/zhanglp0129/lpdrive-server/config"
	"github.com/zhanglp0129/lpdrive-server/model"
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

// MinioNewMultipartUpload 新建分片上传
func MinioNewMultipartUpload(sha256 string) (uploadId string, err error) {
	return MC.NewMultipartUpload(context.Background(),
		config.C.Minio.BucketName,
		sha256,
		minio.PutObjectOptions{})
}

// MinioMultipartUpload 分片上传
func MinioMultipartUpload(multipartUpload *model.MultipartUpload, uploadId string, partId int64, content []byte) error {
	_, err := MC.PutObjectPart(context.Background(),
		config.C.Minio.BucketName, multipartUpload.Sha256,
		uploadId, int(partId), bytes.NewBuffer(content),
		int64(len(content)), minio.PutObjectPartOptions{})
	return err
}

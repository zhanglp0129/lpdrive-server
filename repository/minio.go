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

// MinioObjectExists minio判断对象是否存在。
// 分别返回：是否存在、对象大小、可能发生的异常
func MinioObjectExists(sha256 string) (bool, int64, error) {
	info, err := MC.StatObject(context.Background(), config.C.Minio.BucketName,
		sha256, minio.StatObjectOptions{})
	if minio.ToErrorResponse(err).Code == "NoSuchKey" {
		return false, 0, nil
	} else if err != nil {
		return false, 0, err
	}
	return true, info.Size, nil
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
		uploadId, int(partId)+1, bytes.NewBuffer(content),
		int64(len(content)), minio.PutObjectPartOptions{})
	return err
}

// MinioCompleteUpload minio完成上传
func MinioCompleteUpload(sha256, uploadId string, parts int64) error {
	// 获取ETag
	res, err := MC.ListObjectParts(context.Background(), config.C.Minio.BucketName,
		sha256, uploadId, 0, int(parts))
	if err != nil {
		return err
	}
	minioParts := make([]minio.CompletePart, parts)
	for i := 0; i < int(parts); i++ {
		minioParts[i].PartNumber = res.ObjectParts[i].PartNumber
		minioParts[i].ETag = res.ObjectParts[i].ETag
	}
	_, err = MC.CompleteMultipartUpload(context.Background(), config.C.Minio.BucketName,
		sha256, uploadId, minioParts, minio.PutObjectOptions{})
	return err
}

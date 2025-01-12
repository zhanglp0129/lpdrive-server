package repository

import (
	"context"
	"crypto/sha256"
	"encoding"
	"encoding/json"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	"github.com/zhanglp0129/lpdrive-server/common/constant/minioconstant"
	portaldto "github.com/zhanglp0129/lpdrive-server/dto/portal"
	"github.com/zhanglp0129/lpdrive-server/logger"
	"github.com/zhanglp0129/lpdrive-server/model"
	"time"
)

// RedisSetMultipartUpload 往redis中存入分片上传状态
func RedisSetMultipartUpload(uploadId string, dto portaldto.FilePrepareUploadDTO) error {
	// 新建一个哈希对象
	hasher := sha256.New()
	hashState, err := hasher.(encoding.BinaryMarshaler).MarshalBinary()
	if err != nil {
		return err
	}
	// 将以上数据保存到redis，key为multipart-upload:uploadId，value为 model.MultipartUpload
	key := "multipart-upload:" + uploadId
	value := model.MultipartUpload{
		HashState: hashState,
		Size:      dto.Size,
		Sha256:    dto.Sha256,
		Filename:  dto.Filename,
		PartSize:  dto.PartSize,
		UserID:    dto.UserID,
		DirID:     dto.DirID,
	}
	logger.L.WithField("MultipartUpload", value).Info()
	// redis类型为字符串，json
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	res := RDB.SetNX(context.Background(), key, string(v), minioconstant.UploadExpireDays*24*time.Hour)
	if res.Err() != nil {
		return res.Err()
	} else if !res.Val() {
		return errorconstant.DuplicateUploadId
	}
	return nil
}

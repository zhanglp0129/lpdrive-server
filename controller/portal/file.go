package portalcontroller

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	"github.com/zhanglp0129/lpdrive-server/common/constant/fileconstant"
	"github.com/zhanglp0129/lpdrive-server/common/constant/minioconstant"
	portaldto "github.com/zhanglp0129/lpdrive-server/dto/portal"
	"github.com/zhanglp0129/lpdrive-server/logger"
	portalservice "github.com/zhanglp0129/lpdrive-server/service/portal"
	"github.com/zhanglp0129/lpdrive-server/utils/fileutil"
	"io"
	"strconv"
	"strings"
)

// FileList 获取目录下所有子文件
func FileList(c *gin.Context) (any, error) {
	// 绑定参数
	var dto portaldto.FileListDTO
	err := c.ShouldBindQuery(&dto)
	if err != nil {
		return nil, err
	}
	// 获取用户id
	dto.UserID = c.Value("id").(int64)
	// 调整排序规则
	dto.OrderBy = fileutil.BuildOrderBy(dto.OrderBy, dto.Desc)
	logger.L.WithField("FileListDTO", dto).Info()

	return portalservice.FileList(dto)
}

// FileCreateDirectory 创建目录
func FileCreateDirectory(c *gin.Context) (any, error) {
	// 绑定参数
	var dto portaldto.FileCreateDirectoryEmptyDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		return nil, err
	}
	// 校验文件名
	if err = fileutil.CheckFilename(dto.Name); err != nil {
		return nil, err
	}
	// 获取用户id
	dto.UserID = c.Value("id").(int64)
	logger.L.WithField("FileCreateDirectoryDTO", dto).Info()

	return portalservice.FileCreateDirectory(dto)
}

// FileCreateEmpty 创建空文件
func FileCreateEmpty(c *gin.Context) (any, error) {
	// 绑定参数
	var dto portaldto.FileCreateDirectoryEmptyDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		return nil, err
	}
	// 校验文件名
	if err = fileutil.CheckFilename(dto.Name); err != nil {
		return nil, err
	}
	// 获取用户id
	dto.UserID = c.Value("id").(int64)
	logger.L.WithField("FileCreateEmptyDTO", dto).Info()

	return portalservice.FileCreateEmpty(dto)
}

// FileGetById 根据id获取文件
func FileGetById(c *gin.Context) (any, error) {
	// 获取文件id
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	// 获取用户id
	userId := c.Value("id").(int64)
	return portalservice.FileGetById(id, userId)
}

// FileGetTree 获取文件树
func FileGetTree(c *gin.Context) (any, error) {
	// 获取文件id
	idString, ok := c.GetQuery("id")
	if !ok {
		return nil, errorconstant.IllegalArgument
	}
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return nil, err
	}
	// 获取用户id
	userId := c.Value("id").(int64)
	logger.L.WithField("fileId", id).Info()
	return portalservice.FileGetTree(id, userId)
}

// FileGetPath 获取文件路径
func FileGetPath(c *gin.Context) (any, error) {
	// 获取文件id
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	// 获取用户id
	userId := c.Value("id").(int64)
	return portalservice.FileGetPath(id, userId)
}

// FileGetByPath 根据路径获取文件
func FileGetByPath(c *gin.Context) (any, error) {
	// 获取参数
	path := c.Query("path")
	filenames := strings.Split(path, "/")
	logger.L.WithField("path", path).WithField("filenames", filenames).Info()
	// 获取用户id
	userId := c.Value("id").(int64)
	return portalservice.FileGetByPath(filenames, userId)
}

// FileSearch 搜索文件
func FileSearch(c *gin.Context) (any, error) {
	var dto portaldto.FileSearchDTO
	err := c.ShouldBindQuery(&dto)
	if err != nil {
		return nil, err
	}
	// 获取用户id
	userId := c.Value("id").(int64)
	dto.UserID = userId
	// 调整排序规则
	dto.OrderBy = fileutil.BuildOrderBy(dto.OrderBy, dto.Desc)
	logger.L.WithField("FileSearchDTO", dto).Info()
	return portalservice.FileSearch(dto)
}

// FileSmallUpload 小文件上传
func FileSmallUpload(c *gin.Context) (any, error) {
	var dto portaldto.FileSmallUploadDTO
	err := c.ShouldBind(&dto)
	if err != nil {
		return nil, err
	}
	// 校验文件大小
	if dto.File.Size > fileconstant.SmallFileLimit {
		return nil, errorconstant.FileSizeExceedLimit
	}
	// 检查文件名
	err = fileutil.CheckFilename(dto.File.Filename)
	if err != nil {
		return nil, err
	}
	// 获取用户id
	userId := c.Value("id").(int64)
	dto.UserID = userId

	// 计算文件哈希
	file, err := dto.File.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return nil, err
	}
	hash := hasher.Sum(nil)
	dto.Sha256 = fmt.Sprintf("%x", hash)

	logger.L.WithField("FileSmallUploadDTO", dto).Info()
	err = portalservice.FileSmallUpload(dto)
	return nil, err
}

// FileSmallDownload 小文件下载
func FileSmallDownload(c *gin.Context) (any, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	// 获取用户id
	userId := c.Value("id").(int64)
	logger.L.WithField("id", id).
		WithField("userId", userId).Info()
	return portalservice.FileSmallDownload(id, userId)
}

// FilePrepareUpload 文件预上传
func FilePrepareUpload(c *gin.Context) (any, error) {
	// 获取参数
	var dto portaldto.FilePrepareUploadDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		return nil, err
	}
	// 校验文件名
	err = fileutil.CheckFilename(dto.Filename)
	if err != nil {
		return nil, err
	}
	// 校验分片大小
	if dto.PartSize > minioconstant.MaxPartSize || dto.PartSize < minioconstant.MinPartSize {
		return nil, errorconstant.IllegalPartSize
	}
	// 获取用户id
	userId := c.Value("id").(int64)
	dto.UserID = userId
	logger.L.WithField("FilePrepareUploadDTO", dto).Info()

	return portalservice.FilePrepareUpload(dto)
}

// FileMultipartUpload 文件分片上传
func FileMultipartUpload(c *gin.Context) (any, error) {
	// 获取第几个分片
	partId, err := strconv.ParseInt(c.GetHeader("Part"), 10, 64)
	if err != nil {
		return nil, err
	}
	// 获取upload id
	uploadId := c.Param("uploadId")
	// 获取请求体内容
	content, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(content))
	// 获取用户id
	userId := c.Value("id").(int64)
	logger.L.WithFields(logrus.Fields{
		"partId":        partId,
		"uploadId":      uploadId,
		"contentLength": len(content),
		"userId":        userId,
	}).Info()

	err = portalservice.FileMultipartUpload(partId, uploadId, content, userId)
	return nil, err
}

// FileGetUploadInfo 获取已上传信息
func FileGetUploadInfo(c *gin.Context) (any, error) {
	// 获取upload id
	uploadId := c.Param("uploadId")
	// 获取用户id
	userId := c.Value("id").(int64)
	logger.L.WithFields(logrus.Fields{
		"uploadId": uploadId,
		"userId":   userId,
	}).Info()

	return portalservice.FileGetUploadInfo(uploadId, userId)
}

// FileCompleteUpload 完成分片上传
func FileCompleteUpload(c *gin.Context) (any, error) {
	// 获取upload id
	uploadId := c.Param("uploadId")
	// 获取用户id
	userId := c.Value("id").(int64)
	logger.L.WithFields(logrus.Fields{
		"uploadId": uploadId,
		"userId":   userId,
	}).Info()

	err := portalservice.FileCompleteUpload(uploadId, userId)
	return nil, err
}

package portalcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	portaldto "github.com/zhanglp0129/lpdrive-server/dto/portal"
	"github.com/zhanglp0129/lpdrive-server/logger"
	portalservice "github.com/zhanglp0129/lpdrive-server/service/portal"
	"github.com/zhanglp0129/lpdrive-server/utils/fileutil"
	"strconv"
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
	if dto.OrderBy == "filename" {
		dto.OrderBy = "filename_gbk"
	}
	if dto.Desc {
		dto.OrderBy += " desc"
	}
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

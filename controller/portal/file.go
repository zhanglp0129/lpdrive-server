package portalcontroller

import (
	"github.com/gin-gonic/gin"
	portaldto "github.com/zhanglp0129/lpdrive-server/dto/portal"
	"github.com/zhanglp0129/lpdrive-server/logger"
	portalservice "github.com/zhanglp0129/lpdrive-server/service/portal"
	"github.com/zhanglp0129/lpdrive-server/utils/fileutil"
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

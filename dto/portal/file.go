package portaldto

import "github.com/zhanglp0129/lpdrive-server/dto"

type FileListDTO struct {
	UserID int64
	ID     *int64 `form:"id"`
	dto.PageDTO
	OrderBy string `form:"orderBy" binding:"required,oneof=filename updated_at size"`
}

// FileCreateDirectoryEmptyDTO 创建目录和空文件的参数
type FileCreateDirectoryEmptyDTO struct {
	UserID int64
	DirID  int64  `json:"dirId" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

type FileSearchDTO struct {
	UserID int64
	Name   *string `form:"name"`
	dto.PageDTO
	OrderBy string `form:"orderBy" binding:"required,oneof=filename updated_at size"`
}

package portaldto

import (
	"github.com/zhanglp0129/lpdrive-server/dto"
	"mime/multipart"
)

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

type FileSmallUploadDTO struct {
	UserID   int64
	Sha256   string
	DirID    int64                 `form:"dirId" binding:"required"`
	File     *multipart.FileHeader `form:"file" binding:"required"`
	MimeType string                `form:"mimeType" binding:"required"`
}

type FilePrepareUploadDTO struct {
	UserID   int64
	DirID    int64  `json:"dirId" binging:"required"`
	Filename string `json:"filename" binging:"required"`
	Size     int64  `json:"size" binging:"required"`
	PartSize int64  `json:"partSize" binging:"required"`
	Sha256   string `json:"sha256" binging:"required"`
	MimeType string `json:"mimeType" binding:"required"`
}

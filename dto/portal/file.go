package portaldto

type FileListDTO struct {
	UserID   int64
	ID       *int64 `form:"id"`
	PageNum  int    `form:"pageNum" binding:"required,min=1"`
	PageSize int    `form:"pageSize" binding:"required,min=0"`
	OrderBy  string `form:"orderBy" binding:"required,oneof=filename updated_at size"`
	Desc     bool   `form:"desc"`
}

type FileCreateDirectoryDTO struct {
	UserID int64
	DirID  int64  `json:"dirId" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

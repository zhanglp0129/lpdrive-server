package dto

type PageDTO struct {
	PageNum  int  `form:"pageNum" binding:"required,min=1"`
	PageSize int  `form:"pageSize" binding:"required,min=0"`
	Desc     bool `form:"desc"`
}

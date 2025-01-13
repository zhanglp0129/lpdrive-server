package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanglp0129/lpdrive-server/vo"
	"net/http"
)

// Error 出现错误
func Error(c *gin.Context, err error) {
	_ = c.Error(err)
	c.JSON(http.StatusBadRequest, vo.Error(err.Error()))
	c.Abort()
}

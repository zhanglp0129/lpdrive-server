package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanglp0129/lpdrive-server/vo"
	"net/http"
)

// HandleFunc 自定义路由处理函数
type HandleFunc func(*gin.Context) (any, error)

// 将自定义路由处理函数转化为gin的处理函数
func (h HandleFunc) toGinHandleFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h(c)
		if err != nil {
			// 出现了错误
			_ = c.Error(err)
			c.JSON(http.StatusBadRequest, vo.Error(err.Error()))
			c.Abort()
			return
		}

		// 成功处理
		c.JSON(http.StatusOK, vo.Success(res))
	}
}

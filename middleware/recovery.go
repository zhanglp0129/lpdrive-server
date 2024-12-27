package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanglp0129/lpdrive-server/logger"
	"net/http"
)

// RecoveryMiddleware 崩溃恢复中间件
func RecoveryMiddleware(c *gin.Context) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		// 打印日志
		logger.L.WithField("panic", err).Error()
		c.Status(http.StatusInternalServerError)
		c.Writer.WriteHeaderNow()
		c.Abort()
	}()
}

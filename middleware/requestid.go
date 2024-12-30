package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanglp0129/lpdrive-server/logger"
)

// RequestIdMiddleware 前端可以在请求头加上Request-Id字段，后端会原封不动地放到响应头
func RequestIdMiddleware(c *gin.Context) {
	// 获取requestId
	requestId := c.GetHeader("Request-Id")
	logger.L.WithField("requestId", requestId).Info()
	c.Next()
	// 将requestId放到响应头
	c.Header("Request-Id", requestId)
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zhanglp0129/lpdrive-server/logger"
	"time"
)

func LoggerMiddleware(c *gin.Context) {
	// 获取请求时间
	startTime := time.Now()

	// 处理请求
	c.Next()

	// 计算处理时间
	duration := time.Since(startTime)
	// 获取响应状态码
	code := c.Writer.Status()

	fields := logrus.Fields{
		"method":    c.Request.Method,
		"path":      c.Request.URL.Path,
		"code":      code,
		"client_ip": c.ClientIP(),
		"duration":  duration,
	}

	if code >= 400 {
		// 打印error日志
		logger.L.WithFields(fields).WithField("error", c.Errors.Errors()).Error()
	} else {
		// 打印info日志
		logger.L.WithFields(fields).Info()
	}
}

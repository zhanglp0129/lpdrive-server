package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanglp0129/lpdrive-server/middleware"
)

var (
	R *gin.Engine
)

func init() {
	R = gin.New()
	// 替换默认日志
	R.Use(middleware.LoggerMiddleware)
	// 替换默认崩溃恢复逻辑
	R.Use(middleware.RecoveryMiddleware)
}

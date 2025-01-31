package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanglp0129/lpdrive-server/middleware"
)

var (
	// R 根路由
	R *gin.Engine
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	R = gin.New()
	// 替换默认日志
	R.Use(middleware.LoggerMiddleware)
	// 使用request id中间件
	R.Use(middleware.RequestIdMiddleware)
	// 替换默认崩溃恢复逻辑
	R.Use(middleware.RecoveryMiddleware)

	// 处理用户端路由
	handlePortal()
}

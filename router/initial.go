package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanglp0129/lpdrive-server/middleware"
)

var (
	// R 根路由
	R *gin.Engine
	// 用户端路由
	portal Router
	// 管理端路由
	admin Router
)

func init() {
	R = gin.New()
	// 替换默认日志
	R.Use(middleware.LoggerMiddleware)
	// 替换默认崩溃恢复逻辑
	R.Use(middleware.RecoveryMiddleware)
	// 初始化用户端和管理端路由
	portal.RouterGroup = R.Group("/portal")
	admin.RouterGroup = R.Group("/admin")

	// 为用户端路由绑定中间件
	portal.Use(middleware.LoginMiddleware)
}
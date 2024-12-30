package router

import (
	portalcontroller "github.com/zhanglp0129/lpdrive-server/controller/portal"
	"github.com/zhanglp0129/lpdrive-server/middleware"
)

// 用户端路由
var portal Router

// 处理用户端路由
func handlePortal() {
	// 初始化用户端路由
	portal.RouterGroup = R.Group("/portal")
	// 为用户端路由绑定中间件
	portal.Use(middleware.LoginMiddleware)

	// 处理用户端接口
	handlePortalUser()
}

// 处理用户相关接口
func handlePortalUser() {
	var user Router
	user.RouterGroup = portal.Group("/user")
	user.Post("/login", portalcontroller.UserLogin)
	user.Get("", portalcontroller.UserInfo)
	user.Patch("/password", portalcontroller.UserChangePassword)
	user.Patch("/nickname", portalcontroller.UserChangeNickname)
}

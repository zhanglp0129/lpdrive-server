package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	"github.com/zhanglp0129/lpdrive-server/utils/jwtutil"
)

// LoginMiddleware 用户端登录鉴权中间件
func LoginMiddleware(c *gin.Context) {
	// 排除登录接口
	if c.Request.URL.Path == "/portal/v1/user/login" {
		return
	}

	// 判断是否存在token
	token := c.GetHeader("Authorization")
	if token == "" {
		loginError(c)
		return
	}
	// 解析token
	claims, err := jwtutil.ParseJwtToken(token)
	if err != nil {
		loginError(c)
		return
	}

	// 判断jwt中是否存在id
	idInterface, ok := claims["id"]
	if !ok {
		loginError(c)
		return
	}
	// 判断id是否为int64类型
	id, ok := idInterface.(int64)
	if !ok {
		loginError(c)
		return
	}

	// 将用户id写入上下文
	c.Set("id", id)
}

// 登录错误
func loginError(c *gin.Context) {
	c.String(401, string(errorconstant.LoginTokenError))
	c.Abort()
}

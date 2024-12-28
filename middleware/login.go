package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	"github.com/zhanglp0129/lpdrive-server/logger"
	"github.com/zhanglp0129/lpdrive-server/utils/jwtutil"
	"strconv"
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
		loginError(c, nil, "登录token不存在")
		return
	}
	logger.L.WithField("token", token).Info()

	// 解析token
	claims, err := jwtutil.ParseJwtToken(token)
	if err != nil {
		loginError(c, err, "解析登录token失败")
		return
	}
	logger.L.WithField("claims", claims).Info()

	// 判断id是否存在
	idInterface := claims["id"]
	if idInterface == nil {
		loginError(c, nil, "登录jwt中不存在id")
		return
	}

	// 从jwt中获取用户id
	var id int64
	if idString, ok := idInterface.(string); ok {
		// id是字符串类型
		id, err = strconv.ParseInt(idString, 10, 64)
		if err != nil {
			loginError(c, err, "登录jwt中id字符串解析错误")
			return
		}
	} else if idFloat64, ok := idInterface.(float64); ok {
		// id是数值型
		id = int64(idFloat64)
	} else {
		// id类型错误
		loginError(c, nil, "登录jwt中id类型错误")
		return
	}

	// 将用户id转为int64类型并写入上下文
	c.Set("id", id)
}

// 登录错误
func loginError(c *gin.Context, err error, msg string) {
	// 记录日志
	l := logger.L
	if err != nil {
		l.WithError(err)
	}
	l.Error(msg)

	c.String(401, string(errorconstant.LoginTokenError))
	c.Abort()
}

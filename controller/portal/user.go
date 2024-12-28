package portalcontroller

import (
	"github.com/gin-gonic/gin"
	portaldto "github.com/zhanglp0129/lpdrive-server/dto/portal"
	"github.com/zhanglp0129/lpdrive-server/logger"
	portalservice "github.com/zhanglp0129/lpdrive-server/service/portal"
)

// UserLogin 用户登录接口
func UserLogin(c *gin.Context) (any, error) {
	// 获取请求参数
	var dto portaldto.UserLoginDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		return nil, err
	}
	logger.L.WithField("UserLoginDTO", dto).Info()

	// 调用service登录
	token, err := portalservice.UserLogin(dto.Username, dto.Password)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// UserInfo 查询用户信息接口
func UserInfo(c *gin.Context) (any, error) {
	id := c.Value("id").(int64)
	vo, err := portalservice.UserInfo(id)
	if err != nil {
		return nil, err
	}
	logger.L.WithField("UserInfoVO", vo).Info()
	return vo, nil
}

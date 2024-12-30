package portalcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	portaldto "github.com/zhanglp0129/lpdrive-server/dto/portal"
	"github.com/zhanglp0129/lpdrive-server/logger"
	portalservice "github.com/zhanglp0129/lpdrive-server/service/portal"
	"github.com/zhanglp0129/lpdrive-server/utils/secureutil"
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

// UserChangePassword 修改密码
func UserChangePassword(c *gin.Context) (any, error) {
	// 获取参数
	var dto portaldto.UserChangePasswordDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		return nil, err
	}
	// 校验新密码是否合法
	err = secureutil.CheckPassword(dto.Password)
	if err != nil {
		return nil, err
	}
	// 获取用户id
	id := c.Value("id").(int64)
	dto.ID = id
	logger.L.WithField("UserChangePasswordDTO", dto).Info()

	// 修改密码
	err = portalservice.UserChangePassword(dto)
	return nil, err
}

// UserChangeNickname 修改昵称
func UserChangeNickname(c *gin.Context) (any, error) {
	var dto portaldto.UserChangeNicknameDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		return nil, err
	}
	// 判断昵称长度是否超出上限
	var length int
	for range dto.Nickname {
		length++
		if length > 10 {
			return nil, errorconstant.NicknameLengthExceedLimit
		}
	}
	// 获取用户id
	id := c.Value("id").(int64)
	dto.ID = id
	logger.L.WithField("UserChangeNicknameDTO", dto).Info()

	return nil, portalservice.UserChangeNickname(dto)
}

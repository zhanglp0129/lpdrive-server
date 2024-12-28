package portalservice

import (
	"errors"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	"github.com/zhanglp0129/lpdrive-server/model"
	"github.com/zhanglp0129/lpdrive-server/repository"
	"github.com/zhanglp0129/lpdrive-server/utils/jwtutil"
	"github.com/zhanglp0129/lpdrive-server/utils/secureutil"
	portalvo "github.com/zhanglp0129/lpdrive-server/vo/portal"
	"gorm.io/gorm"
)

func UserLogin(username, password string) (string, error) {
	// 查询数据库，获取密码和盐值
	var user model.User
	err := repository.DB.Select("id", "password", "salt").
		Where("username = ?", username).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 用户名或密码错误
		return "", errorconstant.UsernamePasswordError
	} else if err != nil {
		return "", err
	}

	// 判断加密后的密码是否正确
	encryptPassword := secureutil.EncryptPassword(password, user.Salt)
	if encryptPassword != user.Password {
		return "", errorconstant.UsernamePasswordError
	}

	// 签发jwt token
	claims := map[string]any{
		"id": user.ID,
	}
	return jwtutil.CreateJwtToken(claims)
}

func UserInfo(id int64) (*portalvo.UserInfoVO, error) {
	// 查询数据
	var vo portalvo.UserInfoVO
	err := repository.DB.Model(&model.User{}).
		Where("id = ?", id).Take(&vo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorconstant.UserNotFound
	} else if err != nil {
		return nil, err
	}

	return &vo, nil
}

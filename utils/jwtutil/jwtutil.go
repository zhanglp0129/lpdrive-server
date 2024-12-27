package jwtutil

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	"github.com/zhanglp0129/lpdrive-server/config"
	"time"
)

// CreateJwtToken 创建jwt令牌，传入jwt令牌中存储的数据。
func CreateJwtToken(data map[string]any) (string, error) {
	jwtKey, expireSeconds := config.C.Login.JwtKey, config.C.Login.ExpireSeconds
	claims := jwt.MapClaims{}
	if data != nil {
		for key, value := range data {
			claims[key] = value
		}
	}
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(time.Duration(expireSeconds) * time.Second))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtKey))
}

// ParseJwtToken 解析jwt令牌，并校验
func ParseJwtToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.C.Login.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errorconstant.JwtParseError
	}

	return claims, nil
}

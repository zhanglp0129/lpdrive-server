package secureutil

import (
	"crypto/sha256"
	"fmt"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	"golang.org/x/exp/rand"
)

// GenerateSalt 生成10位的加密盐值
func GenerateSalt() string {
	const start, end = '!', '~'

	chars := make([]byte, 10)
	for i := range chars {
		chars[i] = byte(rand.Intn(end-start+1) + start)
	}
	return string(chars)
}

// EncryptPassword 对密码加密，生成sha256值
func EncryptPassword(password, salt string) string {
	sum := sha256.Sum256([]byte(password + salt))
	return fmt.Sprintf("%x", sum)
}

// CheckPassword 校验密码是否合法。长度至少为8，不能包含非法字符
func CheckPassword(password string) error {
	length := 0
	for _, ch := range password {
		length++
		if !(ch >= '!' && ch <= '~') {
			return errorconstant.IllegalPassword
		}
	}
	if length < 8 {
		return errorconstant.IllegalPassword
	}
	return nil
}

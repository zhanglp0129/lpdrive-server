package secureutil

import (
	"crypto/sha256"
	"fmt"
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

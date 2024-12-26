package model

type Link struct {
	Model
	FileID     int64 `gorm:"not null;index"`
	File       *File
	Name       string `gorm:"type:char(50);not null;uniqueIndex;comment:直链名称，为一个256位36进制的字符串"`
	ExpireTime int64  `gorm:"not null;comment:过期时间，为毫秒级时间戳"`
}

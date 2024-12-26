package model

type User struct {
	Model
	Username    string  `gorm:"uniqueIndex;type:varchar(30);not null;comment:用户名"`
	Password    string  `gorm:"type:char(64);not null;comment:sha256加密后密码"`
	Salt        string  `gorm:"type:char(10);not null;comment:密码加密盐值"`
	NickName    *string `gorm:"type:varchar(10);comment:昵称"`
	Capacity    int64   `gorm:"default:0;not null;comment:总容量，单位为字节"`
	UseCapacity int64   `gorm:"default:0;not null;comment:已使用容量，单位为字节"`
}

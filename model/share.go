package model

type Share struct {
	Model
	FileID int64 `gorm:"uniqueIndex:idx_share,priority:1;comment:分享的文件根路径id"`
	File   *File
	UserID int64 `gorm:"uniqueIndex:idx_share,priority:2;comment:被分享的用户id，需要根据黑白名单确认该用户是否可见"`
	User   *User
}

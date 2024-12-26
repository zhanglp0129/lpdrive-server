package model

type File struct {
	Model
	UserID              int64   `gorm:"index;not null;comment:文件从属用户id"`
	Filename            string  `gorm:"type:varchar(255);unique:uni_file;not null;comment:文件名称"`
	FilenameGBK         []byte  `gorm:"type:varbinary(511);index;not null;comment:GBK编码下文件名称"`
	ObjectName          *string `gorm:"type:char(64);comment:文件在minio中的对象名，为文件sha256值，目录为null"`
	MimeType            *string `gorm:"type:varchar(63);index;comment:文件mine类型，目录为null"`
	Size                int64   `gorm:"index;not null;comment:文件大小，单位为字节。目录大小为0"`
	IsDir               bool    `gorm:"not null;comment:是否为目录"`
	DirID               int64   `gorm:"unique:uni_file;index;not null;comment:文件所处父目录id，根目录为自身"`
	Dir                 *File   // 添加外键：dir_id -> file.id
	OriginalDirID       int64   `gorm:"comment:原父目录id，如果原父目录已删除或该文件未删除，则为null"`
	OriginalDir         *File   // 添加外键：original_dir_id -> file.id
	OriginalFilename    *string `gorm:"index;comment:原文件名称，未删除则为null"`
	OriginalFilenameGBK []byte  `gorm:"comment:GBK编码下的原文件名，未删除则为null"`
	IsShare             bool    `gorm:"default:0;not null;comment:是否被分享"`
	IsWhiteList         *bool   `gorm:"comment:是否启用分享白名单，未分享文件则为null"`
}

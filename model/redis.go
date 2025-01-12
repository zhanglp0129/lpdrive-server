package model

type MultipartUpload struct {
	// 哈希内部状态
	HashState []byte
	// 已上传分片数
	Parts int64
	// 总大小
	Size int64
	// 每个分片大小
	PartSize int64
	// 文件哈希
	Sha256 string
	// 文件名
	Filename string
	// 用户id
	UserID int64
	// 父目录id
	DirID int64
}

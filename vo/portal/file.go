package portalvo

type FileInfo struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"userId"`
	Filename    string  `json:"filename"`
	MimeType    *string `json:"mimeType,omitempty"`
	Size        int64   `json:"size"`
	Sha256      *string `json:"sha256,omitempty"`
	IsDir       bool    `json:"isDir"`
	DirID       int64   `json:"dirId"`
	IsShare     bool    `json:"isShare"`
	IsWhiteList bool    `json:"isWhiteList"`
}

type FileListVO struct {
	DirID int64      `json:"dirId"`
	Items []FileInfo `json:"items"`
}

// FileCreateDirectoryEmptyVO 创建目录和空文件的结果
type FileCreateDirectoryEmptyVO struct {
	ID       int64  `json:"id"`
	SaveName string `json:"saveName"`
}

// FileTreeNode 文件树结点
type FileTreeNode struct {
	FileInfo
	Children []FileTreeNode `json:"children,omitempty" gorm:"-"`
}

type FilePathVO struct {
	ID       int64  `json:"id"`
	Filename string `json:"filename"`
	DirID    *int64 `json:"-"`
}

type FileSearchVO struct {
	Items []FileInfo `json:"items"`
}

type FileSmallDownloadVO struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
	MimeType string `json:"mimeType"`
	Sha256   string `json:"sha256"`
	Size     int64  `json:"size"`
}

type FilePrepareUploadVO struct {
	UploadId *string `json:"uploadId"`
}

type FileGetUploadInfoVO struct {
	Parts    int64  `json:"parts"`
	PartSize int64  `json:"partSize"`
	Sha256   string `json:"sha256"`
	Size     int64  `json:"size"`
}

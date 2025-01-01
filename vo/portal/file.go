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
	Items []FileInfo `json:"items"`
}

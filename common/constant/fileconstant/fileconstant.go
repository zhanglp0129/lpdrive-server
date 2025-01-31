package fileconstant

const (
	// RecycleBin 回收站
	RecycleBin = ":recycle.bin"
	// SmallFileLimit 小文件上限，为20MB
	SmallFileLimit = 20 << 20
)

var (
	// EmptySha256 空文件sha256
	EmptySha256     = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	DefaultMimeType = "text/plain"
)

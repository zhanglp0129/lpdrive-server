package errorconstant

type LpdriveError string

func (l LpdriveError) Error() string {
	return string(l)
}

const (
	FilenameLengthExceedLimit LpdriveError = "文件名长度超出上限"
	IllegalFilename           LpdriveError = "文件名非法"
)

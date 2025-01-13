package errorconstant

type LpdriveError string

func (l LpdriveError) Error() string {
	return string(l)
}

const (
	FilenameLengthExceedLimit LpdriveError = "文件名长度超出上限"
	IllegalFilename           LpdriveError = "文件名非法"
	JwtParseError             LpdriveError = "jwt令牌解析失败"
	LoginTokenError           LpdriveError = "登录信息错误或未登录"
	UsernamePasswordError     LpdriveError = "用户名或密码错误"
	UserNotFound              LpdriveError = "用户不存在"
	IllegalPassword           LpdriveError = "密码非法"
	OldPasswordError          LpdriveError = "旧密码错误"
	NicknameLengthExceedLimit LpdriveError = "昵称长度超出上限"
	TooManyDuplicateNameFiles LpdriveError = "太多重名文件"
	FileNotFound              LpdriveError = "文件不存在"
	IllegalArgument           LpdriveError = "参数非法"
	QueryTimeout              LpdriveError = "查询超时"
	FileSizeExceedLimit       LpdriveError = "文件大小超出上限"
	InsufficientCapacity      LpdriveError = "容量不足"
	DuplicateUploadId         LpdriveError = "上传id重复"
	IllegalPartSize           LpdriveError = "分片大小非法"
	MultipartUploadNotExists  LpdriveError = "分片上传不存在"
	SkipPartsError            LpdriveError = "跳过分片错误"
	UploadFileError           LpdriveError = "上传文件错误"
	FileSizeError             LpdriveError = "文件大小错误"
	RangeError                LpdriveError = "分片错误"
)

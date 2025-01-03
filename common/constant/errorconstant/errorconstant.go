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
)

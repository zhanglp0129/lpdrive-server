package fileutil

import (
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	"strings"
)

// CheckFilename 检查文件名称是否正确，检测通过则返回nil
func CheckFilename(filename string) error {
	// 判断长度是否合法
	if len(filename) == 0 {
		return errorconstant.FilenameLengthExceedLimit
	}
	length := 0
	for range filename {
		length++
		if length > 255 {
			return errorconstant.FilenameLengthExceedLimit
		}
	}

	// 检查是否存在非法字符
	const illegalChars = "<>/\\|:*?"
	if strings.ContainsAny(filename, illegalChars) {
		return errorconstant.IllegalFilename
	}

	// 检查文件名开头和结尾是否存在空格
	if filename[0] == ' ' || filename[len(filename)-1] == ' ' {
		return errorconstant.IllegalFilename
	}
	return nil
}

func BuildOrderBy(orderBy string, desc bool) string {
	res := ""
	// 目录永远在文件之前
	res += "is_dir"
	if desc {
		res += " asc"
	} else {
		res += " desc"
	}
	res += ", "
	// 真正使用排序字段
	if orderBy == "filename" {
		res += "filename_gbk"
	}
	if desc {
		res += " desc"
	}
	return res
}

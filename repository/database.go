package repository

import (
	"fmt"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	"github.com/zhanglp0129/lpdrive-server/utils/dbutil"
	"github.com/zhanglp0129/lpdrive-server/utils/gbkutil"
	"strings"
)

// AttemptAddFile 尝试添加文件。
// attempts为尝试次数，filename为原始文件名，frontExt为添加的序号是否在后缀前面,
// fn为执行的添加文件函数，success为添加成功的回调。
// 返回添加成功的文件名和可能发生的异常
func AttemptAddFile(attempts int, filename string, frontExt bool,
	fn func(name string, gbkName []byte) error) (string, error) {
	// 获取文件名长度
	length := 0
	for range filename {
		length++
	}

	// 尝试添加文件
	for i := 0; i < attempts; i++ {
		name := filename
		if i > 0 {
			// 在文件名上加序号
			num := fmt.Sprintf("(%d)", i)
			// 校验文件名长度
			if length+len(num) > 255 {
				return "", errorconstant.FilenameLengthExceedLimit
			}
			// 拼接文件名
			if frontExt {
				pos := strings.LastIndex(name, ".")
				if pos == -1 {
					pos = len(name)
				}
				name = name[:pos] + num + name[pos:]
			} else {
				name += num
			}
		}
		// 获取文件名的gbk编码
		gbkName, err := gbkutil.StrToGbk(name)
		if err != nil {
			return "", err
		}

		// 执行添加文件函数
		err = fn(name, gbkName)
		if dbutil.IsDuplicateKeyError(err) {
			continue
		} else if err != nil {
			return "", err
		} else {
			return name, nil
		}
	}
	return "", errorconstant.TooManyDuplicateNameFiles
}

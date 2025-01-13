package repository

import (
	"fmt"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	"github.com/zhanglp0129/lpdrive-server/model"
	"github.com/zhanglp0129/lpdrive-server/utils/dbutil"
	"github.com/zhanglp0129/lpdrive-server/utils/gbkutil"
	"gorm.io/gorm"
	"strings"
)

// FileCreateRootDirectory 创建根目录。返回根目录id和可能出现的异常
func FileCreateRootDirectory(tx *gorm.DB, userId int64) (int64, error) {
	// 创建数据模型
	file := model.File{
		UserID:      userId,
		FilenameGBK: make([]byte, 0),
		IsDir:       true,
	}
	// 插入数据
	err := tx.Create(&file).Error
	if err != nil {
		return 0, err
	}
	return file.ID, nil
}

// FileCheckUser 检查文件是否属于该用户
func FileCheckUser(tx *gorm.DB, userId, id int64) error {
	var file model.File
	err := tx.Select("id").Where("user_id = ? and id = ?", userId, id).
		Take(&file).Error
	if err != nil {
		return errorconstant.FileNotFound
	}
	return nil
}

// CheckCapacity 检查用户容量是否足够
func CheckCapacity(tx *gorm.DB, userId, delta int64) error {
	var user model.User
	err := tx.Select("id").Where("id = ? and use_capacity + ? <= capacity", userId, delta).
		Take(&user).Error
	if err != nil {
		return errorconstant.InsufficientCapacity
	}
	return nil
}

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

// DatabaseCreateFile 数据库创建文件。tx必须开启事务
// 该函数会校验用户和容量是否合法，而且会修改使用容量
// 分别返回文件id、保存的文件名、可能发生的异常
func DatabaseCreateFile(tx *gorm.DB, userId int64, sha256 *string, filename string,
	mimeType *string, size int64, isDir bool, dirId int64) (int64, string, error) {
	// 检查父目录是否属于该用户
	err := FileCheckUser(tx, userId, dirId)
	if err != nil {
		return 0, "", err
	}
	// 校验容量
	if size > 0 {
		err = CheckCapacity(tx, userId, size)
		if err != nil {
			return 0, "", err
		}
	}
	// 创建数据模型
	file := model.File{
		UserID:     userId,
		ObjectName: sha256,
		MimeType:   mimeType,
		Size:       size,
		IsDir:      isDir,
		DirID:      &dirId,
	}
	// 添加文件到数据库
	name, err := AttemptAddFile(30, filename, !isDir, func(name string, gbkName []byte) error {
		file.Filename = name
		file.FilenameGBK = gbkName
		return tx.Create(&file).Error
	})
	if err != nil {
		return 0, "", err
	}
	// 增加使用容量
	if size > 0 {
		err = tx.Model(&model.User{}).Where("id = ?", userId).
			Update("use_capacity", gorm.Expr("use_capacity + ?", size)).Error
		if err != nil {
			return 0, "", err
		}
	}
	return file.ID, name, nil
}

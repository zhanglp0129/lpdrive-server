package repository

import (
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	"github.com/zhanglp0129/lpdrive-server/model"
	"gorm.io/gorm"
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

package repository

import (
	"github.com/zhanglp0129/lpdrive-server/model"
	"gorm.io/gorm"
)

// FileCreateRootDirectory 创建根目录。返回根目录id和可能出现的异常
func FileCreateRootDirectory(tx *gorm.DB, userId int64) (int64, error) {
	// 获取根目录id
	id, err := W.GenerateId()
	if err != nil {
		return 0, err
	}
	// 创建数据模型
	file := model.File{
		UserID:      userId,
		FilenameGBK: make([]byte, 0),
		IsDir:       true,
		DirID:       id,
	}
	file.ID = id
	// 插入数据
	err = tx.Create(&file).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

package portalservice

import (
	"errors"
	portaldto "github.com/zhanglp0129/lpdrive-server/dto/portal"
	"github.com/zhanglp0129/lpdrive-server/model"
	"github.com/zhanglp0129/lpdrive-server/repository"
	portalvo "github.com/zhanglp0129/lpdrive-server/vo/portal"
	"gorm.io/gorm"
)

func FileList(dto portaldto.FileListDTO) (portalvo.FileListVO, error) {
	var vo portalvo.FileListVO
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		// 未指定目录，获取根目录id
		if dto.ID == nil {
			var file model.File
			err := tx.Select("id").Where("user_id = ? and id = dir_id", dto.UserID).
				First(&file).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 不存在根目录，创建根目录
				id, err := repository.FileCreateRootDirectory(tx, dto.UserID)
				if err != nil {
					return err
				}
				dto.ID = &id
			} else if err != nil {
				return err
			} else {
				// 存在根目录
				dto.ID = &file.ID
			}
		}

		// 查询数据
		offset := (dto.PageNum - 1) * dto.PageSize
		err := tx.Model(&model.File{}).Select("*", "object_name as sha256").
			Where("user_id = ? and id = ? and id != dir_id", dto.UserID, dto.ID).
			Limit(dto.PageSize).Offset(offset).Order(dto.OrderBy).Find(&vo.Items).Error
		if err != nil {
			return err
		}
		vo.DirID = *dto.ID
		return nil
	})

	if err != nil {
		return portalvo.FileListVO{}, err
	}
	return vo, nil
}

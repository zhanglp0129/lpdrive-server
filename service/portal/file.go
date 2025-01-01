package portalservice

import (
	portaldto "github.com/zhanglp0129/lpdrive-server/dto/portal"
	"github.com/zhanglp0129/lpdrive-server/model"
	"github.com/zhanglp0129/lpdrive-server/repository"
	portalvo "github.com/zhanglp0129/lpdrive-server/vo/portal"
	"gorm.io/gorm"
)

func FileList(dto portaldto.FileListDTO) (portalvo.FileListVO, error) {
	var vo portalvo.FileListVO
	err := repository.DB.Transaction(func(db *gorm.DB) error {
		tx := db.Model(&model.File{}).Select("*", "object_name as sha256")
		// 获取查询条件
		if dto.ID == nil {
			// 未指定目录，默认为根目录
			subQuery := db.Model(&model.File{}).Select("id").
				Where("user_id = ? and id = dir_id", dto.UserID).Limit(1)
			tx = tx.Where("dir_id = (?)", subQuery)
		} else {
			// 指定目录
			tx = tx.Where("user_id = ? and id = ?", dto.UserID, dto.ID)
		}

		// 分页获取根目录下的文件列表
		offset := (dto.PageNum - 1) * dto.PageSize
		err := tx.Limit(dto.PageSize).Offset(offset).
			Order(dto.OrderBy).Find(&vo.Items).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return portalvo.FileListVO{}, err
	}
	return vo, nil
}

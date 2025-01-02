package portalservice

import (
	"errors"
	"fmt"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	portaldto "github.com/zhanglp0129/lpdrive-server/dto/portal"
	"github.com/zhanglp0129/lpdrive-server/model"
	"github.com/zhanglp0129/lpdrive-server/repository"
	"github.com/zhanglp0129/lpdrive-server/utils/dbutil"
	"github.com/zhanglp0129/lpdrive-server/utils/gbkutil"
	portalvo "github.com/zhanglp0129/lpdrive-server/vo/portal"
	"gorm.io/gorm"
)

func FileList(dto portaldto.FileListDTO) (portalvo.FileListVO, error) {
	var vo portalvo.FileListVO
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		// 未指定目录，获取根目录id
		if dto.ID == nil {
			var file model.File
			err := tx.Select("id").Where("user_id = ? and dir_id is null ", dto.UserID).
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
			Where("user_id = ? and dir_id = ?", dto.UserID, dto.ID).
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

func FileCreateDirectory(dto portaldto.FileCreateDirectoryDTO) (portalvo.FileCreateDirectoryVO, error) {
	// 获取文件名长度
	length := 0
	for range dto.Name {
		length++
	}
	// 创建添加数据模型
	file := model.File{
		UserID: dto.UserID,
		IsDir:  true,
		DirID:  &dto.DirID,
	}
	// 检查父目录是否属于该用户
	err := repository.FileCheckUser(repository.DB, dto.UserID, dto.DirID)
	if err != nil {
		return portalvo.FileCreateDirectoryVO{}, err
	}

	// 尝试创建目录
	for i := 0; i <= 30; i++ {
		name := dto.Name
		if i > 0 {
			// 在文件名上加序号，并判断长度
			num := fmt.Sprintf("(%d)", i)
			// 校验文件名长度
			if length+len(num) > 255 {
				return portalvo.FileCreateDirectoryVO{}, errorconstant.FilenameLengthExceedLimit
			}
			name += num
		}
		// 指定文件名
		file.Filename = name
		file.FilenameGBK, err = gbkutil.StrToGbk(name)
		if err != nil {
			return portalvo.FileCreateDirectoryVO{}, err
		}

		// 添加数据
		err = repository.DB.Create(&file).Error
		if dbutil.IsDuplicateKeyError(err) {
			continue
		} else if err != nil {
			return portalvo.FileCreateDirectoryVO{}, err
		} else {
			// 添加成功
			return portalvo.FileCreateDirectoryVO{
				ID:       file.ID,
				SaveName: file.Filename,
			}, nil
		}
	}
	// 重试次数太多
	return portalvo.FileCreateDirectoryVO{}, errorconstant.TooManyDuplicateNameFiles
}

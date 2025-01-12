package portalservice

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/emirpasic/gods/v2/queues/linkedlistqueue"
	"github.com/gabriel-vasile/mimetype"
	"github.com/zhanglp0129/lpdrive-server/common/constant/errorconstant"
	"github.com/zhanglp0129/lpdrive-server/common/constant/fileconstant"
	portaldto "github.com/zhanglp0129/lpdrive-server/dto/portal"
	"github.com/zhanglp0129/lpdrive-server/model"
	"github.com/zhanglp0129/lpdrive-server/repository"
	"github.com/zhanglp0129/lpdrive-server/utils/dbutil"
	"github.com/zhanglp0129/lpdrive-server/utils/fileutil"
	"github.com/zhanglp0129/lpdrive-server/utils/gbkutil"
	portalvo "github.com/zhanglp0129/lpdrive-server/vo/portal"
	"gorm.io/gorm"
	"io"
	"mime"
	"path/filepath"
	"slices"
	"strings"
	"time"
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

func FileCreateDirectory(dto portaldto.FileCreateDirectoryEmptyDTO) (*portalvo.FileCreateDirectoryEmptyVO, error) {
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
		return nil, err
	}

	// 尝试创建目录
	for i := 0; i <= 30; i++ {
		name := dto.Name
		if i > 0 {
			// 在文件名上加序号，并判断长度
			num := fmt.Sprintf("(%d)", i)
			// 校验文件名长度
			if length+len(num) > 255 {
				return nil, errorconstant.FilenameLengthExceedLimit
			}
			name += num
		}
		// 指定文件名
		file.Filename = name
		file.FilenameGBK, err = gbkutil.StrToGbk(name)
		if err != nil {
			return nil, err
		}

		// 添加数据
		err = repository.DB.Create(&file).Error
		if dbutil.IsDuplicateKeyError(err) {
			continue
		} else if err != nil {
			return nil, err
		} else {
			// 添加成功
			return &portalvo.FileCreateDirectoryEmptyVO{
				ID:       file.ID,
				SaveName: file.Filename,
			}, nil
		}
	}
	// 重试次数太多
	return nil, errorconstant.TooManyDuplicateNameFiles
}

func FileCreateEmpty(dto portaldto.FileCreateDirectoryEmptyDTO) (*portalvo.FileCreateDirectoryEmptyVO, error) {
	// 返回结果
	vo := &portalvo.FileCreateDirectoryEmptyVO{}
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		// 获取文件名长度
		length := 0
		for range dto.Name {
			length++
		}
		// 创建添加数据模型
		file := model.File{
			UserID:     dto.UserID,
			ObjectName: &fileconstant.EmptySha256,
			MimeType:   &fileconstant.DefaultMimeType,
			DirID:      &dto.DirID,
		}
		// 检查父目录是否属于该用户
		err := repository.FileCheckUser(tx, dto.UserID, dto.DirID)
		if err != nil {
			return err
		}

		// 尝试创建空文件夹
		for i := 0; i <= 30; i++ {
			name := dto.Name
			if i > 0 {
				// 在文件名上加序号，并判断长度
				num := fmt.Sprintf("(%d)", i)
				// 校验文件名长度
				if length+len(num) > 255 {
					return errorconstant.FilenameLengthExceedLimit
				}
				// 拼接文件名
				pos := strings.LastIndex(name, ".")
				if pos == -1 {
					pos = len(name)
				}
				name = name[:pos] + num + name[pos:]
			}
			// 指定文件名
			file.Filename = name
			file.FilenameGBK, err = gbkutil.StrToGbk(name)
			if err != nil {
				return err
			}

			// 添加数据
			err = tx.Create(&file).Error
			if dbutil.IsDuplicateKeyError(err) {
				continue
			} else if err != nil {
				return err
			} else {
				// 添加成功
				// 将数据写入minio
				err = repository.PutObject(fileconstant.EmptySha256,
					bytes.NewReader(make([]byte, 0)), 0)
				if err != nil {
					return err
				}
				vo.ID = file.ID
				vo.SaveName = file.Filename
				return nil
			}
		}
		// 重试次数太多
		return errorconstant.TooManyDuplicateNameFiles
	})

	if err != nil {
		return nil, err
	}
	return vo, nil
}

func FileGetById(id int64, userId int64) (*portalvo.FileInfo, error) {
	// 查询数据
	var fileInfo portalvo.FileInfo
	err := repository.DB.Model(&model.File{}).Select("*", "object_name as sha256").
		Where("id = ? and user_id = ?", id, userId).Take(&fileInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorconstant.FileNotFound
	} else if err != nil {
		return nil, err
	}
	return &fileInfo, nil
}

func FileGetTree(id int64, userId int64) (*portalvo.FileTreeNode, error) {
	var vo portalvo.FileTreeNode
	// 设置查询3s时间上限
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repository.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx = tx.Model(&model.File{}).
			Select("*", "object_name as sha256").Session(&gorm.Session{})
		// 先获取根节点
		err := tx.Where("id = ? and user_id = ?", id, userId).Take(&vo).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorconstant.FileNotFound
		} else if err != nil {
			return err
		}
		if !vo.IsDir {
			return nil
		}
		// 采用广度优先搜索。队列中存储父目录id和子结点指针
		type queueElementType struct {
			dirId    int64
			children *[]portalvo.FileTreeNode
		}
		queue := linkedlistqueue.New[queueElementType]()
		queue.Enqueue(queueElementType{id, &vo.Children})
		for !queue.Empty() {
			ele, _ := queue.Dequeue()
			err = tx.Where("dir_id = ? and user_id = ?", ele.dirId, userId).
				Find(&ele.children).Error
			if err != nil {
				return err
			}
			// 将子结点插入队列
			for i := range *ele.children {
				if (*ele.children)[i].IsDir {
					queue.Enqueue(queueElementType{
						dirId:    (*ele.children)[i].ID,
						children: &(*ele.children)[i].Children,
					})
				}
			}
		}
		return nil
	})
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, errorconstant.QueryTimeout
	} else if err != nil {
		return nil, err
	}
	return &vo, nil
}

func FileGetPath(id int64, userId int64) ([]portalvo.FilePathVO, error) {
	vos := make([]portalvo.FilePathVO, 0)
	// 设置超时时间 3s
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := repository.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx = tx.Model(&model.File{}).
			Select("id", "filename", "dir_id").Session(&gorm.Session{})
		for {
			// 查询路径
			var vo portalvo.FilePathVO
			err := tx.Where("id = ? and user_id = ?", id, userId).Take(&vo).Error
			if err != nil {
				return err
			}
			if vo.DirID == nil {
				break
			}
			// 将路径添加到结果尾部
			vos = append(vos, vo)
			// 下一次查询其父目录
			id = *vo.DirID
		}
		// 将结果反转
		slices.Reverse(vos)
		return nil
	})
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, errorconstant.QueryTimeout
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorconstant.FileNotFound
	} else if err != nil {
		return nil, err
	}
	return vos, nil
}

func FileGetByPath(path []string, userId int64) (*portalvo.FileInfo, error) {
	vo := portalvo.FileInfo{}
	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := repository.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx = tx.Model(&model.File{}).Select("id").Session(&gorm.Session{})
		// 先获取根目录id
		var id int64
		err := tx.Where("user_id = ? and dir_id is null", userId).Take(&id).Error
		if err != nil {
			return err
		}
		// 按照路径，逐步获取文件id
		length := 0
		for _, filename := range path {
			if len(filename) == 0 {
				continue
			}
			length++
			// 将文件名转为gbk编码，因为有索引
			filenameGBK, err := gbkutil.StrToGbk(filename)
			if err != nil {
				return err
			}
			// 文件文件名查找
			err = tx.Where("user_id = ? and dir_id = ? and filename_gbk = ?", userId, id, filenameGBK).
				Take(&id).Error
			if err != nil {
				return err
			}
		}
		// 路径长度不能为空
		if length == 0 {
			return errorconstant.FileNotFound
		}
		// 根据id获取文件
		return tx.Select("*", "object_name as sha256").
			Where("id = ? and user_id = ?", id, userId).Take(&vo).Error
	})
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, errorconstant.QueryTimeout
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorconstant.FileNotFound
	} else if err != nil {
		return nil, err
	}
	return &vo, nil
}

func FileSearch(dto portaldto.FileSearchDTO) (*portalvo.FileSearchVO, error) {
	var vo portalvo.FileSearchVO
	tx := repository.DB.Model(&model.File{}).Select("*", "object_name as sha256")
	if dto.Name != nil {
		tx = tx.Where("filename like ?", "%"+*dto.Name+"%")
	}
	offset := (dto.PageNum - 1) * dto.PageSize
	err := tx.Where("user_id = ? and dir_id is not null", dto.UserID).
		Offset(offset).Limit(dto.PageSize).Order(dto.OrderBy).Find(&vo.Items).Error
	if err != nil {
		return nil, err
	}
	return &vo, nil
}

func FileSmallUpload(dto portaldto.FileSmallUploadDTO) error {
	return repository.DB.Transaction(func(tx *gorm.DB) error {
		// 获取文件名长度
		length := 0
		for range dto.File.Filename {
			length++
		}
		// 创建添加数据模型
		file := model.File{
			UserID:     dto.UserID,
			ObjectName: &dto.Sha256,
			Size:       dto.File.Size,
			DirID:      &dto.DirID,
		}
		// 根据后缀获取mime type
		mimeType := mime.TypeByExtension(filepath.Ext(dto.File.Filename))
		if mimeType == "" {
			fileReader, err := dto.File.Open()
			if err != nil {
				return err
			}
			defer fileReader.Close()
			m, err := mimetype.DetectReader(fileReader)
			if err != nil {
				return err
			}
			mimeType = m.String()
		}
		file.MimeType = &mimeType
		// 检查父目录是否属于该用户
		err := repository.FileCheckUser(tx, dto.UserID, dto.DirID)
		if err != nil {
			return err
		}
		// 检查容量是否足够
		err = repository.CheckCapacity(tx, dto.UserID, dto.File.Size)
		if err != nil {
			return err
		}

		// 尝试添加文件到数据库
		for i := 0; i <= 30; i++ {
			name := dto.File.Filename
			if i > 0 {
				// 在文件名上加序号，并判断长度
				num := fmt.Sprintf("(%d)", i)
				// 校验文件名长度
				if length+len(num) > 255 {
					return errorconstant.FilenameLengthExceedLimit
				}
				// 拼接文件名
				pos := strings.LastIndex(name, ".")
				if pos == -1 {
					pos = len(name)
				}
				name = name[:pos] + num + name[pos:]
			}
			// 指定文件名
			file.Filename = name
			file.FilenameGBK, err = gbkutil.StrToGbk(name)
			if err != nil {
				return err
			}

			// 添加数据
			err = tx.Create(&file).Error
			if dbutil.IsDuplicateKeyError(err) {
				continue
			} else if err != nil {
				return err
			} else {
				// 添加成功
				// 增加使用容量
				err = tx.Model(&model.User{}).Where("id = ?", dto.UserID).
					Update("use_capacity", gorm.Expr("use_capacity + ?", dto.File.Size)).Error
				if err != nil {
					return err
				}
				// 将数据写入minio
				fileReader, err := dto.File.Open()
				if err != nil {
					return err
				}
				defer fileReader.Close()
				err = repository.PutObject(dto.Sha256, fileReader, dto.File.Size)
				if err != nil {
					return err
				}
				return nil
			}
		}
		// 重试次数太多
		return errorconstant.TooManyDuplicateNameFiles
	})
}

func FileSmallDownload(id int64, userId int64) (*portalvo.FileSmallDownloadVO, error) {
	// 查询数据库
	var vo portalvo.FileSmallDownloadVO
	err := repository.DB.Model(&model.File{}).
		Where("id = ? and user_id = ? and is_dir = 0", id, userId).
		Select("object_name as sha256", "mime_type", "filename", "size").
		Take(&vo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorconstant.FileNotFound
	} else if err != nil {
		return nil, err
	}
	// 判断文件大小
	if vo.Size > fileconstant.SmallFileLimit {
		return nil, errorconstant.FileSizeExceedLimit
	}
	// 从minio获取文件内容
	reader, err := repository.ReadObject(vo.Sha256)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	vo.Content = base64.StdEncoding.EncodeToString(content)
	return &vo, nil
}

func FilePrepareUpload(dto portaldto.FilePrepareUploadDTO) (*portalvo.FilePrepareUploadVO, error) {
	// 校验用户
	err := repository.FileCheckUser(repository.DB, dto.UserID, dto.DirID)
	if err != nil {
		return nil, err
	}
	// 校验容量
	err = repository.CheckCapacity(repository.DB, dto.UserID, dto.Size)
	if err != nil {
		return nil, err
	}

	// 新建一个分片上传
	uploadId, err := repository.MinioNewMultipartUpload(dto.Sha256)
	if err != nil {
		return nil, err
	}

	// 写入redis
	err = repository.RedisSetMultipartUpload(uploadId, dto)
	if err != nil {
		return nil, err
	}
	return &portalvo.FilePrepareUploadVO{UploadId: &uploadId}, nil
}

func FileMultipartUpload(partId int64, uploadId string, content []byte, userId int64) error {
	// 获取redis数据模型
	multipartUpload, hasher, err := repository.RedisGetMultipartUpload(uploadId)
	if err != nil {
		return err
	}
	// 校验分片id
	if multipartUpload.Parts < partId {
		// 跳过分片
		return errorconstant.SkipPartsError
	} else if multipartUpload.Parts > partId {
		// 重复上传
		return nil
	}
	// 校验用户id
	if multipartUpload.UserID != userId {
		return errorconstant.UserNotFound
	}
	// 校验分片大小
	if int64(len(content)) != multipartUpload.PartSize &&
		!fileutil.IsLastPart(partId, multipartUpload.Size, multipartUpload.PartSize, int64(len(content))) {
		return errorconstant.IllegalPartSize
	}

	// 写入minio，minio的part id从1开始
	err = repository.MinioMultipartUpload(multipartUpload, uploadId, partId+1, content)
	if err != nil {
		return err
	}

	// 写回redis
	return repository.RedisUpdateMultipartUpload(multipartUpload, uploadId, hasher, content)
}

package repository

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/zhanglp0129/lpdrive-server/config"
	"github.com/zhanglp0129/lpdrive-server/logger"
	"github.com/zhanglp0129/lpdrive-server/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB redis.UniversalClient
	MC  *minio.Core
)

func init() {
	initDatabase()
	initRedis()
	initMinio()
}

// 初始化数据库连接
func initDatabase() {
	// 读取数据库配置
	dsn := config.C.Database.DSN

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: &logger.GormLogger{},
	})
	if err != nil {
		logger.L.WithField("dsn", dsn).WithError(err).Panicln("数据库连接失败")
	}

	// 迁移表结构
	err = db.AutoMigrate(&model.User{}, &model.File{}, &model.Share{}, &model.Link{})
	if err != nil {
		logger.L.WithError(err).Panicln("数据库表迁移失败")
	}

	// 设置全局数据库实例
	DB = db
}

// 初始化redis
func initRedis() {
	redisConfig := config.C.Redis
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    redisConfig.Addrs,
		DB:       redisConfig.DB,
		Username: redisConfig.Username,
		Password: redisConfig.Password,
	})

	// 判断是否连接成功
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		logger.L.WithField("config", redisConfig).WithError(err).Panicln("redis连接失败")
	}
}

// 初始化minio
func initMinio() {
	minioConfig := config.C.Minio
	mc, err := minio.NewCore(minioConfig.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
	})
	if err != nil {
		logger.L.WithField("minioConfig", minioConfig).WithError(err).Panicln("minio连接失败")
	}

	// 判断存储桶是否存在
	exists, err := mc.BucketExists(context.Background(), minioConfig.BucketName)
	if err != nil {
		logger.L.WithField("minioConfig", minioConfig).WithError(err).Panicln("判断存储桶是否存在时出现错误")
	}
	if !exists {
		// 创建存储桶
		err = mc.MakeBucket(context.Background(), minioConfig.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			logger.L.WithError(err).Panicln("创建存储桶失败")
		}
	}

	MC = mc
}

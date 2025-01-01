package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/zhanglp0129/lpdrive-server/config"
	"github.com/zhanglp0129/lpdrive-server/logger"
	"github.com/zhanglp0129/lpdrive-server/model"
	"github.com/zhanglp0129/snowflake"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	DB  *gorm.DB
	RDB redis.UniversalClient
	W   snowflake.WorkerInterface
)

func init() {
	initDatabase()
	initSnowflake()
	initRedis()
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

// 初始化雪花算法
func initSnowflake() {
	// 设置起始时间
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.L.WithError(err).Panicln("时区加载失败")
	}
	startTime, err := time.ParseInLocation(
		"2006-01-02 15:04:05", "2025-01-01 00:00:00", loc)
	if err != nil {
		logger.L.WithError(err).Panicln("雪花算法起始时间解析失败")
	}

	// 加载雪花算法配置
	c := snowflake.SnowFlakeConfig{
		StartTimestamp: startTime.UnixMilli(),
		TimestampBits:  45,
		MachineIdBits:  8,
		SeqBits:        10,
	}

	// 创建工作结点
	W, err = snowflake.NewWorker(c, 0)
	if err != nil {
		logger.L.WithError(err).Panicln("创建雪花算法工作结点失败")
	}
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

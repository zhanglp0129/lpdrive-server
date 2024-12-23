package repository

import (
	"github.com/zhanglp0129/lpdrive-server/config"
	"github.com/zhanglp0129/lpdrive-server/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	initDatabase()
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
		logger.L.WithError(err).Panicln("数据库连接失败")
	}

	// 设置全局数据库实例
	DB = db
}

package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// C 项目配置实例
var C *Config

// 配置文件路径
const configFile = "./config.yml"

func init() {
	// 读取配置文件
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("配置文件读取失败")
	}

	// 从环境变量读取配置
	viper.SetEnvPrefix("LPDRIVE")
	viper.AutomaticEnv()

	// 设置默认配置
	defaultConfig()

	// 解析配置，生成配置实例
	C = new(Config)
	if err := viper.Unmarshal(C); err != nil {
		panic(err)
	}
}

// 默认配置
func defaultConfig() {
	viper.SetDefault("server.ip", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("login.expire_seconds", 86400)
	viper.SetDefault("database.slow_threshold", 200)
	viper.SetDefault("log.level", "warn")
	viper.SetDefault("log.filename", "logs/lpdrive.log")
	viper.SetDefault("log.max_size", 10)
	viper.SetDefault("log.max_backups", 20)
	viper.SetDefault("log.max_age", 30)
}

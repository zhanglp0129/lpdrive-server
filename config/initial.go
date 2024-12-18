package config

import (
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
		panic(err)
	}

	// 从环境变量读取配置
	viper.SetEnvPrefix("LPDRIVE")
	viper.AutomaticEnv()

	// 解析配置，生成配置实例
	C = new(Config)
	if err := viper.Unmarshal(C); err != nil {
		panic(err)
	}
}

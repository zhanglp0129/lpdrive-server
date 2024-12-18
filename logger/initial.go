package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/zhanglp0129/lpdrive-server/config"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var L *logrus.Logger

// 日志级别
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

func init() {
	// 读取日志相关配置
	logConfig := &config.C.Log

	// 创建 Lumberjack 配置
	logFile := &lumberjack.Logger{
		Filename:   logConfig.Filename,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,
		Compress:   true,
	}

	// 初始化 logrus
	L = logrus.New()
	L.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 输出到文件和标准库
	output := io.MultiWriter(logFile, os.Stdout)
	L.SetOutput(output)

	// 设置日志级别
	var level logrus.Level
	switch logConfig.Level {
	case DebugLevel:
		level = logrus.DebugLevel
	case InfoLevel:
		level = logrus.InfoLevel
	case WarnLevel:
		level = logrus.WarnLevel
	case ErrorLevel:
		level = logrus.ErrorLevel
	default:
		panic("日志级别错误")
	}
	L.SetLevel(level)

	// 显示调用
	L.SetReportCaller(true)
}

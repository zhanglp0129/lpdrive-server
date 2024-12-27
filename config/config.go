package config

// Config 项目配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Admin    AdminConfig    `mapstructure:"admin"`
	Login    LoginConfig    `mapstructure:"login"`
	Database DatabaseConfig `mapstructure:"database"`
	Minio    MinioConfig    `mapstructure:"minio"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	IP   string `mapstructure:"ip"`
	Port uint16 `mapstructure:"port"`
}

// AdminConfig 管理员用户配置
type AdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// LoginConfig 登录配置
type LoginConfig struct {
	JwtKey        string `mapstructure:"jwt_key"`
	ExpireSeconds int64  `mapstructure:"expire_seconds"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN           string `mapstructure:"dsn"`
	SlowThreshold int64  `mapstructure:"slow_threshold"`
}

// MinioConfig minio配置
type MinioConfig struct {
	Endpoint   string `mapstructure:"endpoint"`
	BucketName string `mapstructure:"bucket_name"`
	AccessKey  string `mapstructure:"access_key"`
	SecretKey  string `mapstructure:"secret_key"`
}

// RedisConfig redis相关配置
type RedisConfig struct {
	Addrs    []string `mapstructure:"addrs"`
	DB       int      `mapstructure:"db"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	// 服务器配置
	Server struct {
		// 服务器端口
		Port int `mapstructure:"port" default:"8010"`
	} `mapstructure:"server"`
	// 日志配置
	Logger struct {
		// 日志级别
		Level string `mapstructure:"level" default:"info"`
		// 日志文件路径
		Path string `mapstructure:"path" default:"./log/klog.log"`
		// 日志最大大小（MB）
		MaxSize int `mapstructure:"max_size" default:"100"`
		// 日志最大保留数量
		MaxBackups int `mapstructure:"max_backups" default:"3"`
		// 日志最大保留时间（天）
		MaxAge int `mapstructure:"max_age" default:"30"`
	} `mapstructure:"logger"`
	// 媒体配置
	Media struct {
		// 媒体目录
		MediaDir string `mapstructure:"media_dir" default:"./uploads"`
		// 最大文件大小（MB）
		MaxFileSize int `mapstructure:"max_file_size_mb" default:"10"`
	} `mapstructure:"media"`
	// 数据库配置
	Database struct {
		// 数据库类型
		Type string `mapstructure:"type"`
		// 数据库URL
		Url string `mapstructure:"url"`
	} `mapstructure:"database"`
	// Redis配置
	Redis struct {
		// Redis地址
		Addr string `mapstructure:"addr"`
		// Redis密码
		Password string `mapstructure:"password"`
	} `mapstructure:"redis"`
	// JWT配置
	Jwt struct {
		// JWT密钥
		Secret string `mapstructure:"secret"`
		// JWT过期时间（小时）
		ExpireHour int `mapstructure:"expire_hour"`
	} `mapstructure:"jwt"`
}

var Cfg *Config

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		panic(err)
	}
}

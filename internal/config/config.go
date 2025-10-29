package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// 服务器配置
	Server struct {
		// 服务器端口
		Port int `mapstructure:"port" default:"8010"`
		// 允许的跨域请求源
		Cors []string `mapstructure:"cors" default:"['localhost']"`
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

func Init() error {
	// 设置默认值
	setDefaults()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(Cfg); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	return nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.port", 8010)
	viper.SetDefault("server.cors", []string{"localhost"})
	// 日志默认配置
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.path", "./log/klog.log")
	viper.SetDefault("logger.max_size", 100)
	viper.SetDefault("logger.max_backups", 3)
	viper.SetDefault("logger.max_age", 30)

	// 媒体默认配置
	viper.SetDefault("media.media_dir", "./uploads")
	viper.SetDefault("media.max_file_size_mb", 10)

	// JWT 默认配置
	viper.SetDefault("jwt.expire_hour", 168) // 7天
}

// validateConfig 验证配置的有效性
func validateConfig(cfg *Config) error {
	// 验证必填项
	if cfg.Database.Url == "" {
		return errors.New("数据库URL不能为空")
	}

	if cfg.Jwt.Secret == "" {
		return errors.New("JWT密钥不能为空")
	}

	if len(cfg.Jwt.Secret) < 32 {
		return errors.New("JWT密钥长度至少为32个字符")
	}

	// 验证数值范围
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("服务器端口号无效: %d", cfg.Server.Port)
	}

	if cfg.Media.MaxFileSize < 1 || cfg.Media.MaxFileSize > 100 {
		return fmt.Errorf("文件大小限制无效: %dMB（应在1-100MB之间）", cfg.Media.MaxFileSize)
	}

	if cfg.Jwt.ExpireHour < 1 || cfg.Jwt.ExpireHour > 8760 {
		return fmt.Errorf("JWT过期时间无效: %d小时（应在1-8760小时之间）", cfg.Jwt.ExpireHour)
	}

	return nil
}

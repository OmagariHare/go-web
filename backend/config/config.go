// package config 负责处理应用程序的配置。
// 它使用Viper库从配置文件和环境变量中读取配置，并提供默认值。
package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config 是应用程序所有配置的根结构体。
// 它包含了各个模块的配置，如服务器、数据库、JWT、Casbin、日志和Redis。
type Config struct {
	App      AppConfig      // 应用特定配置
	Server   ServerConfig   // 服务器相关配置
	Database DatabaseConfig // 数据库连接配置
	JWT      JWTConfig      // JWT认证配置
	Casbin   CasbinConfig   // Casbin权限控制配置
	Log      LogConfig      // 日志记录配置
	Redis    RedisConfig    // Redis缓存配置
}

// AppConfig 存储应用级别的配置。
type AppConfig struct {
	DefaultRole string // 新用户注册时的默认角色
}

// ServerConfig 存储服务器相关的配置。
type ServerConfig struct {
	Port           int      // 服务器监听的端口
	AllowedOrigins []string // CORS（跨源资源共享）允许的源列表
}

// DatabaseConfig 存储数据库连接信息。
type DatabaseConfig struct {
	Host     string // 数据库主机地址
	Port     int    // 数据库端口
	User     string // 数据库用户名
	Password string // 数据库密码
	DBName   string // 数据库名称
	SSLMode  string // SSL模式（例如 "disable", "require"）
}

// JWTConfig 存储JWT（JSON Web Token）相关的配置。
type JWTConfig struct {
	Secret     string // 用于签名JWT的密钥
	Expiration int    // JWT的过期时间（以秒为单位）
}

// CasbinConfig 存储Casbin相关的配置。
type CasbinConfig struct {
	Model string // Casbin模型的路径
}

// LogConfig 存储日志记录相关的配置。
type LogConfig struct {
	Level      string // 日志级别 (e.g., "debug", "info", "warn", "error")
	Filename   string // 日志文件路径
	MaxSize    int    // 日志文件的最大大小（以MB为单位）
	MaxBackups int    // 保留的旧日志文件的最大数量
	MaxAge     int    // 保留的旧日志文件的最大天数
	Compress   bool   // 是否压缩旧的日志文件
}

// RedisConfig 存储Redis连接信息。
type RedisConfig struct {
	Addr     string // Redis服务器地址 (e.g., "localhost:6379")
	Password string // Redis密码
	DB       int    // 要使用的Redis数据库编号
}

// LoadConfig 从文件（config.yaml）中加载配置。
// 它会首先设置一系列默认值，然后尝试从配置文件中读取并覆盖这些默认值。
// 如果配置文件不存在或读取失败，将使用默认值并打印警告信息。
func LoadConfig() *Config {
	// 设置配置文件名、类型和路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config") // 在./config/目录下查找配置文件

	// --- 设置各项配置的默认值 ---
	// 如果配置文件中没有相应的值，将使用这些默认值

	// 应用配置
	viper.SetDefault("app.default_role", "user")

	// 服务器配置
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.allowed_origins", []string{"http://localhost:3000"})

	// 数据库配置
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "goweb")
	viper.SetDefault("database.sslmode", "disable")

	// JWT配置
	viper.SetDefault("jwt.expiration", 86400) // 默认24小时

	// 日志配置
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.filename", "./logs/app.log")
	viper.SetDefault("log.maxSize", 10)
	viper.SetDefault("log.maxBackups", 5)
	viper.SetDefault("log.maxAge", 30)
	viper.SetDefault("log.compress", false)

	// Redis配置
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果读取失败，记录一条警告信息，程序将使用默认配置继续运行
		log.Printf("Warning: unable to read config file: %v, using defaults", err)
	}

	// 将读取到的配置信息反序列化到Config结构体中
	config := &Config{
		App: AppConfig{
			DefaultRole: viper.GetString("app.default_role"),
		},
		Server: ServerConfig{
			Port:           viper.GetInt("server.port"),
			AllowedOrigins: viper.GetStringSlice("server.allowed_origins"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.dbname"),
			SSLMode:  viper.GetString("database.sslmode"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("jwt.secret"),
			Expiration: viper.GetInt("jwt.expiration"),
		},
		Casbin: CasbinConfig{
			Model: viper.GetString("casbin.model"),
		},
		Log: LogConfig{
			Level:      viper.GetString("log.level"),
			Filename:   viper.GetString("log.filename"),
			MaxSize:    viper.GetInt("log.maxSize"),
			MaxBackups: viper.GetInt("log.maxBackups"),
			MaxAge:     viper.GetInt("log.maxAge"),
			Compress:   viper.GetBool("log.compress"),
		},
		Redis: RedisConfig{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
	}

	return config
}

// GetDSN 根据数据库配置生成数据源名称（Data Source Name）。
// DSN是用于连接数据库的字符串。
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		c.Database.Host, c.Database.User, c.Database.Password, c.Database.DBName, c.Database.Port, c.Database.SSLMode)
}

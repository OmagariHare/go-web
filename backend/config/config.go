package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Casbin   CasbinConfig
	Log      LogConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Port          int
	AllowedOrigins []string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	Expiration int
}

type CasbinConfig struct {
	Model string
}

type LogConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// 设置默认值
		viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.allowed_origins", []string{"http://localhost:3000"})
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "goweb")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("jwt.secret", "mySecretKey")
	viper.SetDefault("jwt.expiration", 86400)
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.filename", "./logs/app.log")
	viper.SetDefault("log.maxSize", 10)
	viper.SetDefault("log.maxBackups", 5)
	viper.SetDefault("log.maxAge", 30)
	viper.SetDefault("log.compress", false)
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: unable to read config file: %v, using defaults", err)
	}

	config := &Config{
		Server: ServerConfig{
			Port:          viper.GetInt("server.port"),
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

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		c.Database.Host, c.Database.User, c.Database.Password, c.Database.DBName, c.Database.Port, c.Database.SSLMode)
}
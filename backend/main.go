package main

import (
	"fmt"
	"go-web/cache"
	"go-web/config"
	"go-web/database"
	"go-web/routers"
	"go-web/utils"
	"log"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 安全检查：确保JWT密钥已设置
	if cfg.JWT.Secret == "" {
		log.Fatal("FATAL: JWT secret is not configured. Please set 'jwt.secret' in config.yaml or environment variables.")
	}

	// 初始化日志
	utils.InitLogger()
	defer utils.SyncLogger()

	// 连接数据库
	database.ConnectDB(cfg)

	// 连接Redis
	cache.ConnectRedis(cfg)

	// 设置路由
	r := routers.SetupRouter(cfg)

	// 启动服务器
	port := cfg.Server.Port
	log.Printf("Server starting on port %d", port)
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

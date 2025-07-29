package main

import (
	"context"
	"errors"
	"fmt"
	"go-web/config"
	"go-web/database"
	"go-web/routers"
	"go-web/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 安全检查：确保JWT密钥已设置
	if cfg.JWT.Secret == "" {
		log.Fatal("FATAL: JWT secret is not configured. Please set 'jwt.secret' in config.yaml or environment variables.")
	}

	// 初始化日志
	utils.InitLogger(
		cfg.Log.Level,
		cfg.Log.Filename,
		cfg.Log.MaxSize,
		cfg.Log.MaxBackups,
		cfg.Log.MaxAge,
		cfg.Log.Compress,
	)
	defer utils.SyncLogger()

	// 连接数据库
	database.ConnectDB(cfg)

	// 设置路由
	r := routers.SetupRouter(cfg)

	// 启动服务器
	port := cfg.Server.Port
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		// 服务连接
		log.Printf("Server starting on port %d", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

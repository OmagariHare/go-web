package cache

import (
	"context"
	"go-web/config"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func ConnectRedis(cfg *config.Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
}
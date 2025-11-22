package databases

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

const redisPrefix = "refresh:"

func InitRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis-master:6379"
	}
	redisClient = redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis недоступен:", err)
	}
}

func Save(key string, value interface{}, exp time.Duration) {
	redisClient.Set(ctx, redisPrefix+key, value, exp)
}

func Get(key string) (string, error) {
	return redisClient.Get(ctx, redisPrefix+key).Result()
}
func HealthCheck() error {
	_, err := redisClient.Ping(ctx).Result()
	return err
}

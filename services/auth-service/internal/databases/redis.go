package databases

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
)

const redisPrefix = "refresh:"


func InitRedis(redisAddr string) {
	redisClient = redis.NewClient(&redis.Options{Addr: redisAddr})
	
	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		log.Fatalf("redis tracing init failed: %v", err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis недоступен:", err)
	}
}

func Save(key string, value any, exp time.Duration) {
	redisClient.Set(context.Background(), redisPrefix+key, value, exp)
}

func Get(key string) (string, error) {
	return redisClient.Get(context.Background(), redisPrefix+key).Result()
}
func HealthCheck() error {
	return redisClient.Ping(context.Background()).Err()
}

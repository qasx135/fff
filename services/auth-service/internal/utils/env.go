package utils

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret       string
	RedisAddr       string
	GelfEndpoint    string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
	JaegerEndpoint  string
	ServerPort      string
}

const (
	JWT_SECRET        = "fallnacl-secret-change-in-production"
	REDIS_ADDR        = "redis-master:6379"
	GELF_ENDPOINT     = "graylog:12201"
	ACCESS_TOKEN_EXP  = 15 * time.Minute
	REFRESH_TOKEN_EXP = 24 * 7 * time.Hour
	JAEGER_ENDOPOINT  = "jaeger-service:4317"
	SERVER_PORT       = ":8080"
)

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	return &Config{
		JWTSecret:       getEnv("JWT_SECRET", JWT_SECRET),
		RedisAddr:       getEnv("REDIS_ADDR", REDIS_ADDR),
		GelfEndpoint:    getEnv("GELF_ENDPOINT", GELF_ENDPOINT),
		AccessTokenExp:  getEnvAsDuration("ACCESS_TOKEN_EXP", ACCESS_TOKEN_EXP),
		RefreshTokenExp: getEnvAsDuration("REFRESH_TOKEN_EXP", REFRESH_TOKEN_EXP),
		JaegerEndpoint:  getEnv("JAEGER_ENDPOINT", JAEGER_ENDOPOINT),
		ServerPort:      getEnv("SERVER_PORT", SERVER_PORT),
	}
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if dur, err := time.ParseDuration(value); err == nil {
			return dur
		}
	}
	return defaultValue
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

package main

import (
	"context"
	"fmt"
	"github.com/Graylog2/go-gelf/gelf"
	"log"
	_ "net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	jwtSecret   []byte
	gelfLogger  *gelf.Writer
)

const (
	accessTokenExp  = 15 * time.Minute
	refreshTokenExp = 7 * 24 * time.Hour
	redisPrefix     = "refresh:"
)

func initLogger() {
	var err error
	gelfLogger, err = gelf.NewWriter("graylog:12201")
	if err != nil {
		log.Fatal("Не удалось подключиться к Graylog:", err)
	}
}

func logRequest(c *gin.Context) {
	msg := gelf.Message{
		Version: "1.0",
		Level:   gelf.LOG_INFO,
		Short:   fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
		Extra: map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"ip":     c.ClientIP(),
			"ua":     c.Request.UserAgent(),
		},
	}
	if err := gelfLogger.WriteMessage(&msg); err != nil {
		log.Printf("Ошибка отправки в грэйлог: %v", err)
	}
	c.Next()
}

func initRedis() {
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

func initSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "anime-platform-dev-secret-12345" // ← для демо
	}
	jwtSecret = []byte(secret)
}

func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		c.JSON(400, gin.H{"error": "username required"})
		return
	}

	// Access token (HS256)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(accessTokenExp).Unix(),
		"aud": "anime-platform",
	})
	atStr, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": "token signing failed"})
		return
	}

	// Refresh token (просто случайная строка или тоже JWT — здесь упрощённо как строка)
	rtStr := fmt.Sprintf("rt_%d_%s", time.Now().Unix(), username)

	// Храним в Redis
	ctx := context.Background()
	redisClient.Set(ctx, redisPrefix+username, rtStr, refreshTokenExp)

	c.JSON(200, gin.H{
		"access_token":  atStr,
		"refresh_token": rtStr,
	})
}

func refreshHandler(c *gin.Context) {
	oldRT := c.PostForm("refresh_token")
	username := c.PostForm("username") // добавим явно для простоты
	if oldRT == "" || username == "" {
		c.JSON(400, gin.H{"error": "refresh_token and username required"})
		return
	}

	// Проверка в Redis
	ctx := context.Background()
	storedRT, err := redisClient.Get(ctx, redisPrefix+username).Result()
	if err != nil || storedRT != oldRT {
		c.JSON(401, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	// Выдаём новый access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(accessTokenExp).Unix(),
		"aud": "anime-platform",
	})
	atStr, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": "token signing failed"})
		return
	}

	c.JSON(200, gin.H{"access_token": atStr})
}

// Просто health-checks
func healthz(c *gin.Context) { c.Status(200) }
func ready(c *gin.Context)   { c.Status(200) }
func startup(c *gin.Context) { c.Status(200) }

func main() {
	initLogger()
	initRedis()
	initSecret()

	r := gin.Default()
	r.Use(logRequest)

	r.POST("/login", loginHandler)
	r.POST("/refresh", refreshHandler)

	r.GET("/healthz", healthz)
	r.GET("/ready", ready)
	r.GET("/startup", startup)

	log.Println("Auth service запущен на :8080")
	r.Run(":8080")
}

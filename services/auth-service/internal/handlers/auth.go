package handlers

import (
	"auth-service/internal/databases"
	"auth-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	SignManager *utils.JwtManager
}

func (h *AuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		c.JSON(400, gin.H{"error": "username required"})
		return
	}

	// Access token (HS256)
	accessToken, err := h.SignManager.SignAccessToken(username)
	if err != nil {
		c.JSON(500, gin.H{"error": "token signing failed"})
		return
	}

	// Refresh token (просто случайная строка или тоже JWT — здесь упрощённо как строка)
	refreshToken, err := h.SignManager.SignRefreshToken(username)
	if err != nil {
		c.JSON(500, gin.H{"error": "token signing failed"})
		return
	}
	databases.Save(username, refreshToken, h.SignManager.GetRefreshTokenExp())

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	oldRT := c.PostForm("refresh_token")
	username := c.PostForm("username") // добавим явно для простоты
	if oldRT == "" || username == "" {
		c.JSON(400, gin.H{"error": "refresh_token and username required"})
		return
	}
	storedRT, err := databases.Get(username)
	if err != nil || storedRT != oldRT {
		c.JSON(401, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	// Выдаём новый access token
	accessToken, err := h.SignManager.SignAccessToken(username)
	if err != nil {
		c.JSON(500, gin.H{"error": "token signing failed"})
		return
	}

	c.JSON(200, gin.H{"access_token": accessToken})
}

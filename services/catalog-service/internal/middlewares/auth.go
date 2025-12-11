package middlewares

import (
	"log"
	"net/http"
	"strings"

	"catalog-service/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtManager *utils.JwtManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		userID, err := jwtManager.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("INVALID TOKEN: %s", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("username", userID)x
		c.Next()
	}
}

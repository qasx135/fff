package handlers

import (
	"auth-service/internal/databases"
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthCheck() int {
	if err := databases.HealthCheck(); err != nil {
		return http.StatusServiceUnavailable
	}

	return http.StatusOK
}

// Просто health-checks
func Healthz(c *gin.Context) { c.Status(healthCheck()) }
func Ready(c *gin.Context)   { c.Status(healthCheck()) }
func Startup(c *gin.Context) { c.Status(healthCheck()) }

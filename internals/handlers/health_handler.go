package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tacheraSasi/go-api-starter/pkg/database"
)

// HealthHandler will be used to handle health check requests.
type HealthHandler struct{}

// NewHealthHandler will create a new HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "reason": "database_not_initialized"})
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "reason": "database_unavailable"})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "reason": "database_ping_failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

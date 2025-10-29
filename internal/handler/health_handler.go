package handler

import (
	"klog-backend/internal/cache"
	"klog-backend/internal/utils"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck 健康检查接口
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// 检查数据库连接
	sqlDB, err := h.db.DB()
	dbStatus := "healthy"
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "unhealthy"
	}

	// 检查Redis连接
	redisStatus := "not_configured"
	if cache.RedisClient != nil {
		ctx := c.Request.Context()
		if err := cache.RedisClient.Ping(ctx).Err(); err != nil {
			redisStatus = "unhealthy"
		} else {
			redisStatus = "healthy"
		}
	}

	// 总体状态
	overallStatus := "healthy"
	if dbStatus != "healthy" {
		overallStatus = "unhealthy"
	}

	response := map[string]interface{}{
		"status":    overallStatus,
		"timestamp": time.Now().UTC(),
		"services": map[string]string{
			"database": dbStatus,
			"redis":    redisStatus,
		},
	}

	if overallStatus == "unhealthy" {
		utils.ResponseSuccess(c, http.StatusServiceUnavailable, response)
	} else {
		utils.ResponseSuccess(c, http.StatusOK, response)
	}
}

// Metrics 简单的指标接口
func (h *HealthHandler) Metrics(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 获取数据库连接池统计
	sqlDB, _ := h.db.DB()
	dbStats := sqlDB.Stats()

	metrics := map[string]interface{}{
		"timestamp": time.Now().UTC(),
		"go": map[string]interface{}{
			"version":      runtime.Version(),
			"goroutines":   runtime.NumGoroutine(),
			"memory_alloc": m.Alloc / 1024 / 1024,      // MB
			"memory_total": m.TotalAlloc / 1024 / 1024, // MB
			"memory_sys":   m.Sys / 1024 / 1024,        // MB
			"gc_runs":      m.NumGC,
		},
		"database": map[string]interface{}{
			"max_open_connections": dbStats.MaxOpenConnections,
			"open_connections":     dbStats.OpenConnections,
			"in_use":               dbStats.InUse,
			"idle":                 dbStats.Idle,
		},
	}

	utils.ResponseSuccess(c, http.StatusOK, metrics)
}

// ReadyCheck 就绪检查（用于K8s readiness probe）
func (h *HealthHandler) ReadyCheck(c *gin.Context) {
	// 检查数据库是否就绪
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

// LiveCheck 存活检查（用于K8s liveness probe）
func (h *HealthHandler) LiveCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

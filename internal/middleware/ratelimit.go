package middleware

import (
	"context"
	"klog-backend/internal/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter IP限流器
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter 创建IP限流器
// @r 每秒请求数
// @b 突发请求数
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// AddIP 为IP创建限流器
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter

	return limiter
}

// GetLimiter 获取IP的限流器
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()
	return limiter
}

// GlobalIPRateLimiter 全局IP限流器（导出供main.go使用）
var GlobalIPRateLimiter = NewIPRateLimiter(10, 20) // 每秒10个请求，突发20个

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := GlobalIPRateLimiter.GetLimiter(ip)

		if !limiter.Allow() {
			utils.ResponseError(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CleanupIPRateLimiter 定期清理限流器中的过期IP
// @limiter 限流器
// @interval 清理间隔
// @ctx 上下文，用于优雅退出
func CleanupIPRateLimiter(limiter *IPRateLimiter, interval time.Duration, ctx context.Context) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	utils.SugarLogger.Infof("全局限流器清理协程已启动，清理间隔: %v", interval)

	for {
		select {
		case <-ctx.Done():
			utils.SugarLogger.Info("全局限流器清理协程停止")
			return
		case <-ticker.C:
			limiter.mu.Lock()
			oldCount := len(limiter.ips)
			// 清空所有限流器
			limiter.ips = make(map[string]*rate.Limiter)
			limiter.mu.Unlock()
			utils.SugarLogger.Infof("清理限流器中的IP记录，清理数量: %d", oldCount)
		}
	}
}

// RequestSizeLimit 限制请求体大小中间件
// @maxSize 最大字节数（如 10MB = 10 * 1024 * 1024）
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

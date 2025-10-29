package middleware

import (
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

var limiter = NewIPRateLimiter(10, 20) // 每秒10个请求，突发20个

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := limiter.GetLimiter(ip)

		if !limiter.Allow() {
			utils.ResponseError(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CleanupIPRateLimiter 定期清理限流器中的过期IP
func CleanupIPRateLimiter(limiter *IPRateLimiter, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		limiter.mu.Lock()
		// 清空所有限流器
		limiter.ips = make(map[string]*rate.Limiter)
		limiter.mu.Unlock()
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

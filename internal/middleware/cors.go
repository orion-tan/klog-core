package middleware

import (
	"klog-backend/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的Origin
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := config.Cfg.Server.Cors

		// 动态匹配Origin（符合CORS规范）
		if isAllowedOrigin(origin, allowedOrigins) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			// 允许携带凭证
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else if len(allowedOrigins) > 0 && allowedOrigins[0] == "*" {
			// 支持通配符模式
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// 允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// 预检请求缓存时间（秒）
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// 暴露的响应头
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// isAllowedOrigin 检查Origin是否在允许列表中
func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}
	for _, allowed := range allowedOrigins {
		if allowed == origin {
			return true
		}
		// 支持通配符子域名匹配（例如：*.example.com）
		if strings.HasPrefix(allowed, "*.") {
			domain := strings.TrimPrefix(allowed, "*")
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
	}
	return false
}

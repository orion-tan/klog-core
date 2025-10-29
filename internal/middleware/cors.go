package middleware

import (
	"klog-backend/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的源（生产环境应该配置具体的域名）
		allowedOrigins := config.Cfg.Server.Cors
		c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(allowedOrigins, ", "))

		// 允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// 允许携带凭证
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

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

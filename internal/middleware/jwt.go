package middleware

import (
	"klog-backend/internal/cache"
	"klog-backend/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			utils.ResponseError(c, http.StatusUnauthorized, "UNAUTHORIZED", "请求头中缺少Token")
			c.Abort()
			return
		}

		parts := strings.SplitN(token, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ResponseError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Token格式不正确")
			c.Abort()
			return
		}
		token = parts[1]

		// 检查token是否在黑名单中（已登出）
		isBlacklisted, err := cache.IsInBlacklist(token)
		if err == nil && isBlacklisted {
			utils.ResponseError(c, http.StatusUnauthorized, "TOKEN_EXPIRED", "Token已失效")
			c.Abort()
			return
		}

		claims, err := utils.VerifyToken(token)
		if err != nil {
			utils.ResponseError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Token无效")
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

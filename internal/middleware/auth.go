package middleware

import (
	"klog-backend/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthOptional 可选的JWT认证中间件（用于某些接口既支持已认证又支持未认证访问）
func JWTAuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(token, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}
		token = parts[1]

		claims, err := utils.VerifyToken(token)
		if err != nil {
			c.Next()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// AdminAuth 管理员权限验证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			utils.ResponseError(c, http.StatusUnauthorized, "UNAUTHORIZED", "需要登录")
			c.Abort()
			return
		}

		klogClaims := claims.(*utils.KLogClaims)
		if klogClaims.Role != "admin" {
			utils.ResponseError(c, http.StatusForbidden, "FORBIDDEN", "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}


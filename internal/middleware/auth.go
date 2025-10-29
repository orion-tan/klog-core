package middleware

import (
	"klog-backend/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthOptional 可选的JWT认证中间件
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

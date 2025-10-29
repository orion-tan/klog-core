package middleware

import (
	"klog-backend/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
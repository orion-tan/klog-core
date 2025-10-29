package middleware

import (
	"klog-backend/internal/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CommentRateLimiter 评论专用限流器（防刷）
type CommentRateLimiter struct {
	visitors map[string]*CommentVisitor
	mu       *sync.RWMutex
}

// CommentVisitor 访客信息
type CommentVisitor struct {
	lastCommentTime time.Time
	commentCount    int
	resetTime       time.Time
}

var commentLimiter = &CommentRateLimiter{
	visitors: make(map[string]*CommentVisitor),
	mu:       &sync.RWMutex{},
}

// CommentRateLimitMiddleware 评论限流中间件
// 规则：
// 1. 同一IP每分钟最多1条评论
// 2. 同一IP每小时最多10条评论
func CommentRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		commentLimiter.mu.Lock()
		visitor, exists := commentLimiter.visitors[ip]

		if !exists {
			// 新访客
			commentLimiter.visitors[ip] = &CommentVisitor{
				lastCommentTime: now,
				commentCount:    1,
				resetTime:       now.Add(1 * time.Hour),
			}
			commentLimiter.mu.Unlock()
			c.Next()
			return
		}

		// 检查是否需要重置计数器（每小时重置）
		if now.After(visitor.resetTime) {
			visitor.commentCount = 0
			visitor.resetTime = now.Add(1 * time.Hour)
		}

		// 规则1：距离上次评论不足1分钟
		if now.Sub(visitor.lastCommentTime) < 1*time.Minute {
			commentLimiter.mu.Unlock()
			utils.ResponseError(c, http.StatusTooManyRequests, "COMMENT_TOO_FAST", "评论过于频繁，请1分钟后再试")
			c.Abort()
			return
		}

		// 规则2：1小时内评论超过10条
		if visitor.commentCount >= 10 {
			commentLimiter.mu.Unlock()
			utils.ResponseError(c, http.StatusTooManyRequests, "COMMENT_LIMIT_EXCEEDED", "评论次数过多，请1小时后再试")
			c.Abort()
			return
		}

		// 更新访客信息
		visitor.lastCommentTime = now
		visitor.commentCount++
		commentLimiter.mu.Unlock()

		c.Next()
	}
}

// CleanupCommentLimiter 定期清理过期的访客记录
func CleanupCommentLimiter() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		commentLimiter.mu.Lock()
		for ip, visitor := range commentLimiter.visitors {
			// 清理2小时前的记录
			if now.Sub(visitor.lastCommentTime) > 2*time.Hour {
				delete(commentLimiter.visitors, ip)
			}
		}
		commentLimiter.mu.Unlock()
	}
}

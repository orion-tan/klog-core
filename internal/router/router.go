package router

import (
	"klog-backend/internal/handler"
	"klog-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RouterHandlers 路由处理器集合
type RouterHandlers struct {
	AuthHandler     *handler.AuthHandler
	PostHandler     *handler.PostHandler
	CategoryHandler *handler.CategoryHandler
	TagHandler      *handler.TagHandler
	CommentHandler  *handler.CommentHandler
	MediaHandler    *handler.MediaHandler
	UserHandler     *handler.UserHandler
	HealthHandler   *handler.HealthHandler
	SettingHandler  *handler.SettingHandler
}

func SetupRouter(handlers *RouterHandlers) *gin.Engine {
	router := gin.Default()

	// 全局中间件：CORS跨域支持
	router.Use(middleware.CORSMiddleware())

	// 健康检查和监控端点（不需要认证和限流）
	router.GET("/health", handlers.HealthHandler.HealthCheck)
	router.GET("/health/ready", handlers.HealthHandler.ReadyCheck)
	router.GET("/health/live", handlers.HealthHandler.LiveCheck)
	router.GET("/metrics", handlers.HealthHandler.Metrics)

	// 全局中间件：请求体大小限制（30MB）
	router.Use(middleware.RequestSizeLimit(30 * 1024 * 1024))

	// 全局中间件：限流
	router.Use(middleware.RateLimitMiddleware())

	// api/v1 路由组
	apiV1 := router.Group("/api/v1")
	{
		// 认证路由
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", handlers.AuthHandler.Register) // 仅供首次设置管理员账号
			auth.POST("/login", handlers.AuthHandler.Login)
			auth.GET("/me", middleware.JWTAuth(), handlers.AuthHandler.GetMe)
			auth.POST("/logout", middleware.JWTAuth(), handlers.AuthHandler.Logout)
		}

		// 文章路由
		posts := apiV1.Group("/posts")
		{
			posts.GET("", middleware.JWTAuthOptional(), handlers.PostHandler.GetPosts)
			posts.POST("", middleware.JWTAuth(), handlers.PostHandler.CreatePost)
			posts.GET("/:id", middleware.JWTAuthOptional(), handlers.PostHandler.GetPostByID)
			posts.PUT("/:id", middleware.JWTAuth(), handlers.PostHandler.UpdatePost)
			posts.DELETE("/:id", middleware.JWTAuth(), handlers.PostHandler.DeletePost)

			// 文章评论（使用相同的参数名 :id）
			posts.GET("/:id/comments", handlers.CommentHandler.GetCommentsByPostID)
			posts.POST("/:id/comments", middleware.JWTAuthOptional(), middleware.CommentRateLimitMiddleware(), handlers.CommentHandler.CreateComment)
		}

		// 分类路由
		categories := apiV1.Group("/categories")
		{
			categories.GET("", handlers.CategoryHandler.GetCategories)
			categories.POST("", middleware.JWTAuth(), handlers.CategoryHandler.CreateCategory)
			categories.PUT("/:id", middleware.JWTAuth(), handlers.CategoryHandler.UpdateCategory)
			categories.DELETE("/:id", middleware.JWTAuth(), handlers.CategoryHandler.DeleteCategory)
		}

		// 标签路由
		tags := apiV1.Group("/tags")
		{
			tags.GET("", handlers.TagHandler.GetTags)
			tags.POST("", middleware.JWTAuth(), handlers.TagHandler.CreateTag)
			tags.PUT("/:id", middleware.JWTAuth(), handlers.TagHandler.UpdateTag)
			tags.DELETE("/:id", middleware.JWTAuth(), handlers.TagHandler.DeleteTag)
		}

		// 评论路由
		comments := apiV1.Group("/comments")
		{
			comments.PUT("/:id", middleware.JWTAuth(), handlers.CommentHandler.UpdateCommentStatus)
			comments.DELETE("/:id", middleware.JWTAuth(), handlers.CommentHandler.DeleteComment)
		}

		// 媒体库路由
		media := apiV1.Group("/media")
		{
			// 静态文件服务（上传的文件）
			media.GET("/i/:filename", handlers.MediaHandler.ServeMedia)
			// 文件上传
			media.POST("/upload", middleware.JWTAuth(), handlers.MediaHandler.UploadMedia)
			media.GET("", middleware.JWTAuth(), handlers.MediaHandler.GetMediaList)
			media.DELETE("/:id", middleware.JWTAuth(), handlers.MediaHandler.DeleteMedia)
		}

		// 用户路由
		users := apiV1.Group("/users")
		{
			users.GET("/:id", handlers.UserHandler.GetUserByID)
			users.PUT("/:id", middleware.JWTAuth(), handlers.UserHandler.UpdateUser)
		}

		// 设置路由（仅管理员可访问）
		settings := apiV1.Group("/settings")
		{
			settings.GET("", middleware.JWTAuth(), handlers.SettingHandler.GetAllSettings)
			settings.GET("/:key", middleware.JWTAuth(), handlers.SettingHandler.GetSettingByKey)
			settings.PUT("", middleware.JWTAuth(), handlers.SettingHandler.UpsertSetting)
			settings.PUT("/batch", middleware.JWTAuth(), handlers.SettingHandler.BatchUpsertSettings)
			settings.DELETE("/:key", middleware.JWTAuth(), handlers.SettingHandler.DeleteSetting)
		}
	}

	return router
}

package main

import (
	"context"
	"fmt"
	"klog-backend/internal/cache"
	"klog-backend/internal/config"
	"klog-backend/internal/database"
	"klog-backend/internal/handler"
	"klog-backend/internal/middleware"
	"klog-backend/internal/repository"
	"klog-backend/internal/router"
	"klog-backend/internal/services"
	"klog-backend/internal/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	if err := config.Init(); err != nil {
		panic(fmt.Sprintf("配置初始化失败: %v", err))
	}
	// 初始化日志
	utils.InitLogger()

	// 初始化Redis缓存（可选）
	if err := cache.InitRedis(); err != nil {
		utils.SugarLogger.Warnf("Redis初始化失败，将不使用缓存: %v", err)
	} else if cache.RedisClient != nil {
		utils.SugarLogger.Info("Redis连接成功")
		defer cache.CloseRedis()
	}

	// 启动评论限流清理协程
	go middleware.CleanupCommentLimiter()

	// 初始化数据库
	db, err := gorm.Open(sqlite.Open(config.Cfg.Database.Url), &gorm.Config{})
	if err != nil {
		utils.SugarLogger.Error("初始化数据库失败:", err)
		return
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(db); err != nil {
		utils.SugarLogger.Error("数据库迁移失败:", err)
		return
	}

	// 初始化仓库层
	authRepo := repository.NewAuthRepository(db)
	postRepo := repository.NewPostRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	userRepo := repository.NewUserRepository(db)

	// 初始化服务层
	authService := services.NewAuthService(authRepo)
	postService := services.NewPostService(postRepo, tagRepo, categoryRepo)
	categoryService := services.NewCategoryService(categoryRepo, postRepo)
	tagService := services.NewTagService(tagRepo, postRepo)
	commentService := services.NewCommentService(commentRepo, postRepo)
	mediaService := services.NewMediaService(mediaRepo)
	userService := services.NewUserService(userRepo)

	// 初始化处理器层
	authHandler := handler.NewAuthHandler(authService)
	postHandler := handler.NewPostHandler(postService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	tagHandler := handler.NewTagHandler(tagService)
	commentHandler := handler.NewCommentHandler(commentService)
	mediaHandler := handler.NewMediaHandler(mediaService)
	userHandler := handler.NewUserHandler(userService)
	healthHandler := handler.NewHealthHandler(db)

	// 设置路由
	handlers := &router.RouterHandlers{
		AuthHandler:     authHandler,
		PostHandler:     postHandler,
		CategoryHandler: categoryHandler,
		TagHandler:      tagHandler,
		CommentHandler:  commentHandler,
		MediaHandler:    mediaHandler,
		UserHandler:     userHandler,
		HealthHandler:   healthHandler,
	}

	r := router.SetupRouter(handlers)

	// 创建 HTTP 服务器
	addr := fmt.Sprintf(":%d", config.Cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 在 goroutine 中启动服务器
	go func() {
		utils.SugarLogger.Infof("服务器启动成功，监听地址: %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.SugarLogger.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.SugarLogger.Info("正在关闭服务器...")

	// 设置 5 秒的超时时间用于处理现有请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.SugarLogger.Fatalf("服务器强制关闭: %v", err)
	}

	utils.SugarLogger.Info("服务器已优雅退出")
}

package main

import (
	"fmt"
	"klog-backend/internal/config"
	"klog-backend/internal/database"
	"klog-backend/internal/handler"
	"klog-backend/internal/repository"
	"klog-backend/internal/router"
	"klog-backend/internal/services"
	"klog-backend/internal/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	config.Init()
	// 初始化日志
	utils.InitLogger()
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

	// 设置路由
	handlers := &router.RouterHandlers{
		AuthHandler:     authHandler,
		PostHandler:     postHandler,
		CategoryHandler: categoryHandler,
		TagHandler:      tagHandler,
		CommentHandler:  commentHandler,
		MediaHandler:    mediaHandler,
		UserHandler:     userHandler,
	}

	r := router.SetupRouter(handlers)

	addr := fmt.Sprintf(":%d", config.Cfg.Server.Port)
	utils.SugarLogger.Infof("服务器启动成功，监听地址: %s", addr)
	if err := r.Run(addr); err != nil {
		utils.SugarLogger.Error("服务器启动失败:", err)
		panic(err)
	}
}

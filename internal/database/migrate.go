package database

import (
	"klog-backend/internal/model"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Category{},
		&model.Tag{},
		&model.PostTag{},
		&model.Comment{},
		&model.Setting{},
		&model.Media{},
		&model.UserIdentity{},
	)
}


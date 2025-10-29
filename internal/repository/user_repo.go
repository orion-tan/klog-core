package repository

import (
	"klog-backend/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetUserByID 根据ID获取用户
// @userID 用户ID
// @return 用户, 错误
func (r *UserRepository) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers 获取用户列表（带分页）
// @page 页码
// @limit 每页数量
// @return 用户列表, 总数, 错误
func (r *UserRepository) GetUsers(page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.DB.Model(&model.User{})

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit).Order("created_at DESC")

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser 更新用户信息
// @user 用户
// @return 错误
func (r *UserRepository) UpdateUser(user *model.User) error {
	return r.DB.Save(user).Error
}


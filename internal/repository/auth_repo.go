package repository

import (
	"klog-backend/internal/model"

	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

// CreateUser 创建用户
// @user 新建用户
// @return 新建用户, 错误
func (r *AuthRepository) CreateUser(user *model.User) error {
	return r.DB.Create(user).Error
}

// GetUserByUsername 根据用户名获取用户
// @username 用户名
// @return 用户, 错误
func (r *AuthRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
// @email 邮箱
// @return 用户, 错误
func (r *AuthRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
// @userID 用户ID
// @return 用户, 错误
func (r *AuthRepository) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
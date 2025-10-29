package services

import (
	"errors"
	"klog-backend/internal/api"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetUserByID 根据ID获取用户
// @userID 用户ID
// @return 用户, 错误
func (s *UserService) GetUserByID(userID uint) (*model.User, error) {
	return s.userRepo.GetUserByID(userID)
}

// GetUsers 获取用户列表
// @page 页码
// @limit 每页数量
// @return 用户列表, 总数, 错误
func (s *UserService) GetUsers(page, limit int) ([]model.User, int64, error) {
	return s.userRepo.GetUsers(page, limit)
}

// UpdateUser 更新用户信息
// @userID 用户ID
// @req 更新用户请求
// @return 用户, 错误
func (s *UserService) UpdateUser(userID uint, req *api.UserUpdateRequest) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}


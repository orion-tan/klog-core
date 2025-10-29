package services

import (
	"errors"
	"klog-backend/internal/api"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"
	"klog-backend/internal/utils"
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

// UpdateUser 更新用户信息
// @userID 用户ID
// @req 更新用户请求
// @return 用户, 错误
func (s *UserService) UpdateUser(userID uint, req *api.UserUpdateRequest) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证和更新字段
	if req.Username != nil && *req.Username != "" {
		if !utils.ValidateUsername(*req.Username) {
			return nil, errors.New("用户名格式不正确")
		}
		user.Username = *req.Username
	}
	if req.Nickname != nil && *req.Nickname != "" {
		user.Nickname = *req.Nickname
	}
	if req.AvatarURL != nil {
		if *req.AvatarURL == "" {
			user.AvatarURL = nil
		} else {
			if !utils.ValidateURL(*req.AvatarURL) {
				return nil, errors.New("头像URL格式不正确")
			}
			user.AvatarURL = req.AvatarURL
		}
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.Email != nil && *req.Email != "" {
		if !utils.ValidateEmail(*req.Email) {
			return nil, errors.New("邮箱格式不正确")
		}
		user.Email = *req.Email
	}

	// 检查旧密码是否正确
	if req.OldPassword != nil && *req.OldPassword != "" {
		if !utils.ComparePasswordHash(*req.OldPassword, user.Password) {
			return nil, errors.New("旧密码不正确")
		}
	}
	if req.NewPassword != nil && *req.NewPassword != "" {
		hashedPassword, err := utils.GeneratePasswordHash(*req.NewPassword)
		if err != nil {
			return nil, errors.New("生成密码哈希失败")
		}
		user.Password = hashedPassword
	}
	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

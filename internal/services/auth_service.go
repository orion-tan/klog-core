package services

import (
	"errors"
	"klog-backend/internal/api"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"
	"klog-backend/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

// Register 注册用户
// @req 注册请求
// @return 用户, 错误
func (s *AuthService) Register(req *api.UserRegisterRequest) (user *model.User, err error) {
	// 检查用户是否存在
	_, err = s.authRepo.GetUserByUsername(req.Username)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("用户已存在")
	}
	// 检查邮箱是否存在
	_, err = s.authRepo.GetUserByEmail(req.Email)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("邮箱已存在")
	}

	// 生成密码哈希
	hashedPassword, err := utils.GeneratePasswordHash(req.Password)
	if err != nil {
		return nil, errors.New("生成密码哈希失败")
	}

	// 创建用户
	user = &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
	}

	if err := s.authRepo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login 用户登录
// @req 登录请求
// @return token, 用户, 错误
func (s *AuthService) Login(req *api.UserLoginRequest) (token string, user *model.User, err error) {
	// 尝试通过用户名查找
	user, err = s.authRepo.GetUserByUsername(req.Login)
	if err != nil {
		// 尝试通过邮箱查找
		user, err = s.authRepo.GetUserByEmail(req.Login)
		if err != nil {
			return "", nil, errors.New("用户名或邮箱不存在")
		}
	}

	// 验证密码
	if !utils.ComparePasswordHash(req.Password, user.Password) {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != "active" {
		return "", nil, errors.New("用户已被禁用")
	}

	// 生成Token
	token, err = utils.GenerateToken(user.ID, user.Username, user.Role, user.Status)
	if err != nil {
		return "", nil, errors.New("生成Token失败")
	}

	return token, user, nil
}

// GetUserByID 根据ID获取用户信息
// @userID 用户ID
// @return 用户, 错误
func (s *AuthService) GetUserByID(userID uint) (*model.User, error) {
	return s.authRepo.GetUserByID(userID)
}

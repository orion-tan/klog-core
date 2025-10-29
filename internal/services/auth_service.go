package services

import (
	"errors"
	"klog-backend/internal/api"
	"klog-backend/internal/cache"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"
	"klog-backend/internal/utils"
	"time"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

// Register 注册用户（仅允许首个用户注册）
// @req 注册请求
// @return 用户, 错误
func (s *AuthService) Register(req *api.UserRegisterRequest) (user *model.User, err error) {
	// 检查系统中是否已存在用户（单用户限制）
	count, err := s.authRepo.CountUsers()
	if err != nil {
		return nil, errors.New("数据库查询失败")
	}
	if count > 0 {
		return nil, errors.New("系统已关闭注册")
	}

	// 验证用户名和邮箱
	if !utils.ValidateUsername(req.Username) {
		return nil, errors.New("用户名格式不正确")
	}
	if !utils.ValidateEmail(req.Email) {
		return nil, errors.New("邮箱格式不正确")
	}

	// 生成密码哈希
	hashedPassword, err := utils.GeneratePasswordHash(req.Password)
	if err != nil {
		return nil, errors.New("生成密码哈希失败")
	}

	// 创建唯一管理员用户
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
	if req.Login == "" {
		return "", nil, errors.New("登录信息不能为空")
	}
	if req.Password == "" {
		return "", nil, errors.New("密码不能为空")
	}
	// 验证登录信息
	if !utils.ValidateUsername(req.Login) && !utils.ValidateEmail(req.Login) {
		return "", nil, errors.New("用户名或邮箱格式不正确")
	}
	if !utils.ValidatePassword(req.Password) {
		return "", nil, errors.New("密码格式不正确")
	}

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
		return "", nil, errors.New("密码错误")
	}

	// 生成Token
	token, err = utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, errors.New("生成JWT令牌失败")
	}

	return token, user, nil
}

// GetUserByID 根据ID获取用户信息
// @userID 用户ID
// @return 用户, 错误
func (s *AuthService) GetUserByID(userID uint) (*model.User, error) {
	return s.authRepo.GetUserByID(userID)
}

// Logout 用户登出（将token加入黑名单）
// @token JWT token字符串
// @expiresAt token过期时间
// @return 错误
func (s *AuthService) Logout(token string, expiresAt time.Time) error {
	// 计算token剩余有效期
	duration := time.Until(expiresAt)
	if duration <= 0 {
		// token已过期，无需加入黑名单
		return nil
	}

	// 将token加入黑名单
	return cache.AddToBlacklist(token, duration)
}

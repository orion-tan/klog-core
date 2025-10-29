package handler

import (
	"klog-backend/internal/api"
	"klog-backend/internal/services"
	"klog-backend/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register 用户注册（仅供首次设置管理员账号）
func (h *AuthHandler) Register(c *gin.Context) {
	var req api.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}
	user, err := h.authService.Register(&req)
	if err != nil {
		// 记录注册失败的审计日志
		utils.LogRegister(c, req.Username, false)
		utils.ResponseError(c, http.StatusBadRequest, "REGISTER_FAILED", err.Error())
		return
	}
	// 记录注册成功的审计日志
	utils.LogRegister(c, req.Username, true)
	utils.ResponseSuccess(c, http.StatusCreated, api.UserRegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req api.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}
	token, _, err := h.authService.Login(&req)
	if err != nil {
		// 记录登录失败的审计日志
		utils.LogLogin(c, req.Login, false)
		utils.ResponseError(c, http.StatusUnauthorized, "LOGIN_FAILED", err.Error())
		return
	}
	// 记录登录成功的审计日志
	utils.LogLogin(c, req.Login, true)
	utils.ResponseSuccess(c, http.StatusOK, api.UserLoginResponse{
		Token: token,
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		utils.ResponseError(c, http.StatusUnauthorized, "UNAUTHORIZED", "未找到用户信息")
		return
	}
	klogClaims := claims.(*utils.KLogClaims)

	user, err := h.authService.GetUserByID(klogClaims.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "USER_NOT_FOUND", "用户不存在")
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, api.UserGetMeResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从请求头提取token
	token := c.GetHeader("Authorization")
	if token == "" {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_TOKEN", "缺少token")
		return
	}

	// 提取Bearer token
	parts := strings.SplitN(token, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_TOKEN", "token格式不正确")
		return
	}
	token = parts[1]

	// 获取claims以获取过期时间
	claims, exists := c.Get("claims")
	if !exists {
		utils.ResponseError(c, http.StatusUnauthorized, "UNAUTHORIZED", "未找到用户信息")
		return
	}
	klogClaims := claims.(*utils.KLogClaims)

	// 调用服务层登出
	if err := h.authService.Logout(token, klogClaims.ExpiresAt.Time); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "LOGOUT_FAILED", "登出失败")
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, map[string]string{
		"message": "登出成功",
	})
}

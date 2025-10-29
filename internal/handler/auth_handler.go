package handler

import (
	"klog-backend/internal/api"
	"klog-backend/internal/services"
	"klog-backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req api.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}
	user, err := h.authService.Register(&req)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "REGISTER_FAILED", err.Error())
		return
	}
	utils.ResponseSuccess(c, http.StatusCreated, api.UserRegisterResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
		Status:    user.Status,
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
	utils.ResponseSuccess(c, http.StatusOK, api.UserRegisterResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
		Status:    user.Status,
	})
}

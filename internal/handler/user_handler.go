package handler

import (
	"klog-backend/internal/api"
	"klog-backend/internal/services"
	"klog-backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUserByID 获取单个用户信息
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_ID", "无效的用户ID")
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "USER_NOT_FOUND", "用户不存在")
		return
	}

	// 返回公开信息
	response := api.UserGetMeResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
	}

	utils.ResponseSuccess(c, http.StatusOK, response)
}

// UpdateUser 更新用户信息
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_ID", "无效的用户ID")
		return
	}

	var req api.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	claims, _ := c.Get("claims")
	klogClaims := claims.(*utils.KLogClaims)

	// 检查权限
	if klogClaims.UserID != uint(id) {
		utils.ResponseError(c, http.StatusForbidden, "FORBIDDEN", "无权更新此用户信息")
		return
	}

	user, err := h.userService.UpdateUser(uint(id), &req)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "UPDATE_USER_FAILED", err.Error())
		return
	}

	response := api.UserGetMeResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
	}

	utils.ResponseSuccess(c, http.StatusOK, response)
}

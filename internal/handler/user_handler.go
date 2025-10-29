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

// GetUsers 获取用户列表
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, total, err := h.userService.GetUsers(page, limit)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "GET_USERS_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, api.PaginatedResponse{
		Total: int(total),
		Page:  page,
		Limit: limit,
		Data:  users,
	})
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
	response := api.UserRegisterResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
		Status:    user.Status,
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

	// 检查权限：只能更新自己的信息，或者管理员可以更新任何用户
	if klogClaims.UserID != uint(id) && klogClaims.Role != "admin" {
		utils.ResponseError(c, http.StatusForbidden, "FORBIDDEN", "无权更新此用户信息")
		return
	}

	// 非管理员不能修改状态
	if klogClaims.Role != "admin" {
		req.Status = ""
	}

	user, err := h.userService.UpdateUser(uint(id), &req)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "UPDATE_USER_FAILED", err.Error())
		return
	}

	response := api.UserRegisterResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
		Status:    user.Status,
	}

	utils.ResponseSuccess(c, http.StatusOK, response)
}


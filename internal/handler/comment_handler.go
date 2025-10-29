package handler

import (
	"klog-backend/internal/api"
	"klog-backend/internal/services"
	"klog-backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateComment 创建评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_POST_ID", "无效的文章ID")
		return
	}

	var req api.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	// 获取用户ID（如果已认证）
	var userID *uint
	claims, exists := c.Get("claims")
	if exists {
		klogClaims := claims.(*utils.KLogClaims)
		userID = &klogClaims.UserID
	}

	// 获取IP地址
	ip := c.ClientIP()

	comment, err := h.commentService.CreateComment(uint(postID), &req, userID, ip)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "CREATE_COMMENT_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, comment)
}

// GetCommentsByPostID 获取文章的评论列表
func (h *CommentHandler) GetCommentsByPostID(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_POST_ID", "无效的文章ID")
		return
	}

	comments, err := h.commentService.GetCommentsByPostID(uint(postID))
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "GET_COMMENTS_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, comments)
}

// UpdateCommentStatus 更新评论状态
func (h *CommentHandler) UpdateCommentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_ID", "无效的评论ID")
		return
	}

	var req api.CommentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	comment, err := h.commentService.UpdateCommentStatus(uint(id), &req)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "UPDATE_COMMENT_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, comment)
}

// DeleteComment 删除评论
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_ID", "无效的评论ID")
		return
	}

	if err := h.commentService.DeleteComment(uint(id)); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "DELETE_COMMENT_FAILED", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}


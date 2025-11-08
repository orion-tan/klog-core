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
		if klogClaims, ok := claims.(*utils.KLogClaims); ok {
			userID = &klogClaims.UserID
		}
	}

	// 游客评论验证
	if userID == nil {
		// 游客必须提供姓名和邮箱
		if req.Name == "" {
			utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", "游客评论必须提供姓名")
			return
		}
		if req.Email == "" {
			utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", "游客评论必须提供邮箱")
			return
		}
		// 验证姓名长度
		if len(req.Name) < 2 || len(req.Name) > 50 {
			utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", "姓名长度应在2-50个字符之间")
			return
		}
	}

	// 验证评论内容长度
	if len(req.Content) < 1 || len(req.Content) > 1000 {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", "评论内容长度应在1-1000个字符之间")
		return
	}

	// 验证Markdown内容安全性
	if !utils.ValidateMarkdownContent(req.Content) {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_CONTENT", "评论内容包含不安全的标签或脚本")
		return
	}

	// 清理评论内容（移除危险元素）
	req.Content = utils.SanitizeMarkdownContent(req.Content)

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

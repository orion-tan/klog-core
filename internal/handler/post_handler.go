package handler

import (
	"klog-backend/internal/api"
	"klog-backend/internal/services"
	"klog-backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// CreatePost 创建文章
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req api.PostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	claims, _ := c.Get("claims")
	klogClaims := claims.(*utils.KLogClaims)

	post, err := h.postService.CreatePost(&req, klogClaims.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "CREATE_POST_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, post)
}

// GetPostByID 获取文章详情
func (h *PostHandler) GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_ID", "无效的文章ID")
		return
	}

	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "POST_NOT_FOUND", "文章不存在")
		return
	}

	// 检查是否有权限访问非公开文章
	if post.Status != "published" {
		_, exists := c.Get("claims")
		if !exists {
			utils.ResponseError(c, http.StatusForbidden, "FORBIDDEN", "无权访问此文章")
			return
		}
	}

	utils.ResponseSuccess(c, http.StatusOK, post)
}

// GetPosts 获取文章列表
func (h *PostHandler) GetPosts(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	categorySlug := c.Query("category")
	tagSlug := c.Query("tag")
	sortBy := c.DefaultQuery("sortBy", "published_at")
	order := c.DefaultQuery("order", "desc")
	detail, _ := strconv.Atoi(c.DefaultQuery("detail", "0"))

	// 检查权限
	_, exists := c.Get("claims")
	if !exists {
		// 未认证用户只能看到已发布的文章
		status = "published"
	}

	posts, total, err := h.postService.GetPosts(page, limit, status, categorySlug, tagSlug, sortBy, order)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "GET_POSTS_FAILED", err.Error())
		return
	}

	if detail == 0 {
		for i := range posts {
			posts[i].Content = ""
		}
	}

	utils.ResponseSuccess(c, http.StatusOK, api.PaginatedResponse{
		Total: int(total),
		Page:  page,
		Limit: limit,
		Data:  posts,
	})
}

// UpdatePost 更新文章
func (h *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_ID", "无效的文章ID")
		return
	}

	var req api.PostUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	// 检查文章是否存在
	_, err = h.postService.GetPostByID(uint(id))
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "POST_NOT_FOUND", "文章不存在")
		return
	}

	updatedPost, err := h.postService.UpdatePost(uint(id), &req)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "UPDATE_POST_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, updatedPost)
}

// DeletePost 删除文章
func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_ID", "无效的文章ID")
		return
	}

	if err := h.postService.DeletePost(uint(id)); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "DELETE_POST_FAILED", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

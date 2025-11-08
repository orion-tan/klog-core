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

// GetPosts 获取文章列表（支持游标分页和传统分页）
func (h *PostHandler) GetPosts(c *gin.Context) {
	// 解析通用查询参数
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	categorySlug := c.Query("category")
	tagSlug := c.Query("tag")
	sortBy := c.DefaultQuery("sortBy", "published_at")
	order := c.DefaultQuery("order", "desc")
	detail, _ := strconv.Atoi(c.DefaultQuery("detail", "0"))

	// 限制每页最大数量
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 10
	}

	// 检查权限
	_, exists := c.Get("claims")
	if !exists {
		// 未认证用户只能看到已发布的文章
		status = "published"
	}

	// 判断使用哪种分页模式
	cursor := c.Query("cursor")
	if cursor != "" || c.Query("page") == "" {
		// 使用游标分页
		h.getPostsByCursor(c, cursor, limit, status, categorySlug, tagSlug, sortBy, order, detail)
	} else {
		// 使用offset分页
		h.getPostsByOffset(c, limit, status, categorySlug, tagSlug, sortBy, order, detail)
	}
}

// getPostsByCursor 游标分页获取文章列表
func (h *PostHandler) getPostsByCursor(c *gin.Context, cursor string, limit int, status, categorySlug, tagSlug, sortBy, order string, detail int) {
	posts, nextCursor, hasMore, err := h.postService.GetPostsByCursor(cursor, limit, status, categorySlug, tagSlug, sortBy, order)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "GET_POSTS_FAILED", err.Error())
		return
	}

	// 如果不需要详情，清除content字段
	if detail == 0 {
		for i := range posts {
			posts[i].Content = ""
		}
	}

	utils.ResponseSuccess(c, http.StatusOK, api.CursorPaginatedResponse{
		Data:       posts,
		NextCursor: nextCursor,
		PrevCursor: nil, // 暂不支持反向游标
		HasMore:    hasMore,
		Limit:      limit,
	})
}

// getPostsByOffset 传统offset分页获取文章列表
func (h *PostHandler) getPostsByOffset(c *gin.Context, limit int, status, categorySlug, tagSlug, sortBy, order string, detail int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	posts, total, err := h.postService.GetPosts(page, limit, status, categorySlug, tagSlug, sortBy, order)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "GET_POSTS_FAILED", err.Error())
		return
	}

	// 如果不需要详情，清除content字段
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

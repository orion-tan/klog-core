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
		claims, exists := c.Get("claims")
		if !exists {
			utils.ResponseError(c, http.StatusForbidden, "FORBIDDEN", "无权访问此文章")
			return
		}
		klogClaims := claims.(*utils.KLogClaims)
		// 只有作者或管理员可以访问非公开文章
		if post.AuthorID != klogClaims.UserID && klogClaims.Role != "admin" {
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
	authorIDStr := c.Query("author")
	categorySlug := c.Query("category")
	tagSlug := c.Query("tag")
	sortBy := c.DefaultQuery("sortBy", "published_at")
	order := c.DefaultQuery("order", "desc")
	detail, _ := strconv.Atoi(c.DefaultQuery("detail", "0"))

	// 如果用户已认证，可以查看自己的所有文章
	var authorID *uint
	var uid uint
	if authorIDStr != "" {
		id, err := strconv.ParseUint(authorIDStr, 10, 32)
		if err == nil {
			uid = uint(id)
			authorID = &uid
		}
	}

	// 检查权限
	claims, exists := c.Get("claims")
	if exists {
		klogClaims := claims.(*utils.KLogClaims)
		// 非管理员只能看到已发布的文章或自己的文章
		if klogClaims.Role != "admin" {
			// 普通用户：只能查看已发布的文章，或者查看自己的文章
			if status != "" && status != "published" {
				// 如果查询非发布状态，只能查看自己的
				authorID = &klogClaims.UserID
			} else if authorID != nil && *authorID != klogClaims.UserID {
				// 如果查询其他作者，强制只显示已发布的
				status = "published"
			}
		}
	} else {
		// 未认证用户只能看到已发布的文章
		status = "published"
		authorID = nil
	}

	posts, total, err := h.postService.GetPosts(page, limit, status, categorySlug, tagSlug, sortBy, order, authorID)
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

	// 检查权限
	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "POST_NOT_FOUND", "文章不存在")
		return
	}

	claims, _ := c.Get("claims")
	klogClaims := claims.(*utils.KLogClaims)

	// 只有作者或管理员可以更新文章
	if post.AuthorID != klogClaims.UserID && klogClaims.Role != "admin" {
		utils.ResponseError(c, http.StatusForbidden, "FORBIDDEN", "无权更新此文章")
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

	// 检查权限
	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "POST_NOT_FOUND", "文章不存在")
		return
	}

	claims, _ := c.Get("claims")
	klogClaims := claims.(*utils.KLogClaims)

	// 只有作者或管理员可以删除文章
	if post.AuthorID != klogClaims.UserID && klogClaims.Role != "admin" {
		utils.ResponseError(c, http.StatusForbidden, "FORBIDDEN", "无权删除此文章")
		return
	}

	if err := h.postService.DeletePost(uint(id)); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "DELETE_POST_FAILED", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

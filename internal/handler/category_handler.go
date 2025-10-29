package handler

import (
	"klog-backend/internal/api"
	"klog-backend/internal/services"
	"klog-backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// CreateCategory 创建分类
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req api.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	category, err := h.categoryService.CreateCategory(&req)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "CREATE_CATEGORY_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, category)
}

// GetCategories 获取所有分类
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "GET_CATEGORIES_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, categories)
}

// UpdateCategory 更新分类
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_ID", "无效的分类ID")
		return
	}

	var req api.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	category, err := h.categoryService.UpdateCategory(uint(id), &req)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "UPDATE_CATEGORY_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, category)
}

// DeleteCategory 删除分类
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_ID", "无效的分类ID")
		return
	}

	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "DELETE_CATEGORY_FAILED", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}


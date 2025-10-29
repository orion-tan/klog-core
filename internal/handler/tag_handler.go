package handler

import (
	"klog-backend/internal/api"
	"klog-backend/internal/services"
	"klog-backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagService *services.TagService
}

func NewTagHandler(tagService *services.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

// CreateTag 创建标签
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req api.TagCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	tag, err := h.tagService.CreateTag(&req)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "CREATE_TAG_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, tag)
}

// GetTags 获取所有标签
func (h *TagHandler) GetTags(c *gin.Context) {
	tags, err := h.tagService.GetTags()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "GET_TAGS_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, tags)
}

// UpdateTag 更新标签
func (h *TagHandler) UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_ID", "无效的标签ID")
		return
	}

	var req api.TagUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	tag, err := h.tagService.UpdateTag(uint(id), &req)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "UPDATE_TAG_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, tag)
}

// DeleteTag 删除标签
func (h *TagHandler) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_ID", "无效的标签ID")
		return
	}

	if err := h.tagService.DeleteTag(uint(id)); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "DELETE_TAG_FAILED", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}


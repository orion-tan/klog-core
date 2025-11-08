package handler

import (
	"klog-backend/internal/api"
	"klog-backend/internal/services"
	"klog-backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	settingService *services.SettingService
}

func NewSettingHandler(settingService *services.SettingService) *SettingHandler {
	return &SettingHandler{settingService: settingService}
}

// GetAllSettings 获取所有设置
func (h *SettingHandler) GetAllSettings(c *gin.Context) {
	settings, err := h.settingService.GetAllSettings()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "GET_SETTINGS_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, settings)
}

// GetSettingByKey 根据Key获取设置
func (h *SettingHandler) GetSettingByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_KEY", "设置键不能为空")
		return
	}

	value, err := h.settingService.GetSettingByKey(key)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "SETTING_NOT_FOUND", "设置不存在")
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, gin.H{
		"key":   key,
		"value": value,
	})
}

// UpsertSetting 创建或更新设置
func (h *SettingHandler) UpsertSetting(c *gin.Context) {
	var req api.SettingUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	if err := h.settingService.UpsertSetting(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "UPSERT_SETTING_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, gin.H{
		"key":   req.Key,
		"value": req.Value,
		"type":  req.Type,
	})
}

// BatchUpsertSettings 批量创建或更新设置
func (h *SettingHandler) BatchUpsertSettings(c *gin.Context) {
	var req api.SettingBatchUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	if err := h.settingService.BatchUpsertSettings(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "BATCH_UPSERT_SETTINGS_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, gin.H{
		"count": len(req.Settings),
	})
}

// DeleteSetting 删除设置
func (h *SettingHandler) DeleteSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_KEY", "设置键不能为空")
		return
	}

	if err := h.settingService.DeleteSetting(key); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "DELETE_SETTING_FAILED", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

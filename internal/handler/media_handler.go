package handler

import (
	"encoding/base64"
	"fmt"
	"io"
	"klog-backend/internal/api"
	"klog-backend/internal/services"
	"klog-backend/internal/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	mediaService *services.MediaService
}

func NewMediaHandler(mediaService *services.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: mediaService}
}

// UploadMedia 上传媒体文件
func (h *MediaHandler) UploadMedia(c *gin.Context) {
	contentType := c.ContentType()

	var fileData *services.FileData
	var err error

	if strings.Contains(contentType, "multipart/form-data") {
		// 提取multipart文件数据
		fileData, err = h.extractMultipartFile(c)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "INVALID_FILE", err.Error())
			return
		}
	} else if strings.Contains(contentType, "application/json") {
		// 提取base64文件数据
		fileData, err = h.extractBase64File(c)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "INVALID_FILE", err.Error())
			return
		}
	} else {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_CONTENT_TYPE", "不支持的内容类型")
		return
	}

	// 交给Service层处理
	media, err := h.mediaService.SaveMediaFile(fileData)
	if err != nil {
		// 记录上传失败的审计日志
		utils.LogFileOperation(c, "upload_file", "", fileData.FileName, false)
		utils.ResponseError(c, http.StatusInternalServerError, "UPLOAD_FAILED", err.Error())
		return
	}

	// 记录上传成功的审计日志
	utils.LogFileOperation(c, "upload_file", fmt.Sprintf("%d", media.ID), media.FileName, true)

	// 格式化响应
	response := map[string]interface{}{
		"id":         media.ID,
		"file_name":  media.FileName,
		"file_path":  media.FilePath,
		"file_hash":  media.FileHash,
		"url":        "/media/i/" + media.FilePath,
		"mime_type":  media.MimeType,
		"size":       media.Size,
		"created_at": media.CreatedAt,
	}

	utils.ResponseSuccess(c, http.StatusCreated, response)
}

// extractMultipartFile 从multipart请求中提取文件数据（流式处理）
// @c Gin上下文
// @return 文件数据, 错误
func (h *MediaHandler) extractMultipartFile(c *gin.Context) (*services.FileData, error) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("获取上传文件失败: %v", err)
	}

	// 打开文件获取Reader
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	// 暂时不关闭文件，Service层处理完后再关闭

	// 对于大文件（>1MB），使用流式处理
	const largeFileThreshold = 1 * 1024 * 1024

	if fileHeader.Size > largeFileThreshold {
		// 大文件：使用Reader流式处理
		return &services.FileData{
			FileName: fileHeader.Filename,
			Reader:   file, // 使用流式读取
			MimeType: fileHeader.Header.Get("Content-Type"),
			Size:     fileHeader.Size,
		}, nil
	}

	// 小文件：读取到内存中（用于MIME类型验证）
	fileBytes, err := io.ReadAll(file)
	file.Close()
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %v", err)
	}

	return &services.FileData{
		FileName:  fileHeader.Filename,
		FileBytes: fileBytes,
		MimeType:  fileHeader.Header.Get("Content-Type"),
		Size:      fileHeader.Size,
	}, nil
}

// extractBase64File 从JSON请求中提取base64编码的文件数据
// @c Gin上下文
// @return 文件数据, 错误
func (h *MediaHandler) extractBase64File(c *gin.Context) (*services.FileData, error) {
	var req struct {
		FileName string `json:"file_name" binding:"required"`
		Data     string `json:"data" binding:"required"`
		MimeType string `json:"mime_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, fmt.Errorf("解析请求失败: %v", err)
	}

	// 解码base64数据
	fileBytes, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return nil, fmt.Errorf("解码base64数据失败: %v", err)
	}

	return &services.FileData{
		FileName:  req.FileName,
		FileBytes: fileBytes,
		MimeType:  req.MimeType,
		Size:      int64(len(fileBytes)),
	}, nil
}

// GetMediaList 获取媒体文件列表
func (h *MediaHandler) GetMediaList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// 参数验证
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	mediaList, total, err := h.mediaService.GetMediaList(page, limit)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "GET_MEDIA_LIST_FAILED", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, api.PaginatedResponse{
		Total: int(total),
		Page:  page,
		Limit: limit,
		Data:  mediaList,
	})
}

// DeleteMedia 删除媒体文件
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	// 解析和验证ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "INVALID_ID", "无效的媒体文件ID")
		return
	}

	// 删除媒体文件（包括物理文件和数据库记录）
	if err := h.mediaService.DeleteMediaWithFile(uint(id)); err != nil {
		// 记录删除失败的审计日志
		utils.LogFileOperation(c, "delete_file", idStr, "", false)

		if err == services.ErrMediaNotFound {
			utils.ResponseError(c, http.StatusNotFound, "MEDIA_NOT_FOUND", "媒体文件不存在")
		} else {
			utils.ResponseError(c, http.StatusInternalServerError, "DELETE_MEDIA_FAILED", err.Error())
		}
		return
	}

	// 记录删除成功的审计日志
	utils.LogFileOperation(c, "delete_file", idStr, "", true)

	c.Status(http.StatusNoContent)
}

// ServeMedia 提供媒体文件访问
func (h *MediaHandler) ServeMedia(c *gin.Context) {
	fileName := c.Param("filename")

	// 基本的路径验证
	fileName = filepath.Base(fileName)
	if strings.Contains(fileName, "..") || strings.ContainsAny(fileName, "/\\") {
		c.Status(http.StatusBadRequest)
		return
	}

	// 通过Service层获取文件路径
	filePath, err := h.mediaService.GetMediaFilePath(fileName)
	if err != nil {
		if err == services.ErrMediaNotFound {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusForbidden)
		}
		return
	}

	c.File(filePath)
}

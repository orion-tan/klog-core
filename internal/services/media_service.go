package services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"klog-backend/internal/config"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrMediaNotFound    = errors.New("媒体文件不存在")
	ErrPermissionDenied = errors.New("权限不足")
	ErrInvalidFileType  = errors.New("不支持的文件类型")
	ErrFileTooLarge     = errors.New("文件过大")
	ErrInvalidFileName  = errors.New("无效的文件名")
)

// 允许的MIME类型
var AllowedMimeTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/gif":     true,
	"image/webp":    true,
	"image/svg+xml": true,
}

// 允许的文件扩展名
var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".svg":  true,
}

// FileData 文件数据传输对象
type FileData struct {
	FileName  string
	FileBytes []byte
	MimeType  string
	Size      int64
}

type MediaService struct {
	mediaRepo *repository.MediaRepository
}

func NewMediaService(mediaRepo *repository.MediaRepository) *MediaService {
	return &MediaService{mediaRepo: mediaRepo}
}

// SaveMediaFile 保存媒体文件（包括物理文件和数据库记录）
// @fileData 文件数据
// @uploaderID 上传者ID
// @return 媒体文件, 错误
func (s *MediaService) SaveMediaFile(fileData *FileData, uploaderID uint) (*model.Media, error) {
	// 1. 验证文件
	if err := s.validateFile(fileData); err != nil {
		return nil, err
	}

	// 2. 计算文件哈希
	fileHash, err := s.calculateFileHash(fileData.FileBytes)
	if err != nil {
		return nil, fmt.Errorf("计算文件哈希失败: %v", err)
	}

	// 3. 检查文件是否已存在（去重）
	existingMedia, err := s.mediaRepo.GetMediaByHash(fileHash)
	if err == nil && existingMedia != nil {
		// 文件已存在，直接返回
		return existingMedia, nil
	}

	// 4. 生成唯一文件名
	fileName := s.generateFileName(uploaderID, fileData.FileName)

	// 5. 确保上传目录存在
	uploadDir := config.Cfg.Media.MediaDir
	if err := s.ensureUploadDir(uploadDir); err != nil {
		return nil, err
	}

	// 6. 保存物理文件
	filePath := filepath.Join(uploadDir, fileName)
	if err := s.saveFile(filePath, fileData.FileBytes); err != nil {
		return nil, err
	}

	// 7. 创建数据库记录
	media := &model.Media{
		UploaderID: uploaderID,
		FileName:   fileData.FileName,
		FilePath:   fileName,
		FileHash:   fileHash,
		MimeType:   fileData.MimeType,
		Size:       fileData.Size,
	}

	if err := s.mediaRepo.CreateMedia(media); err != nil {
		// 数据库插入失败，清理已保存的物理文件
		_ = os.Remove(filePath)
		return nil, fmt.Errorf("创建媒体记录失败: %v", err)
	}

	return media, nil
}

// validateFile 验证文件
// @fileData 文件数据
// @return 错误
func (s *MediaService) validateFile(fileData *FileData) error {
	// 检查文件大小
	maxSize := int64(config.Cfg.Media.MaxFileSize * 1024 * 1024)
	if fileData.Size > maxSize {
		return fmt.Errorf("%w: 文件大小超过 %dMB", ErrFileTooLarge, config.Cfg.Media.MaxFileSize)
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(fileData.FileName))
	if !AllowedExtensions[ext] {
		return fmt.Errorf("%w: %s", ErrInvalidFileType, ext)
	}

	// 检查MIME类型
	if !AllowedMimeTypes[fileData.MimeType] {
		return fmt.Errorf("%w: %s", ErrInvalidFileType, fileData.MimeType)
	}

	return nil
}

// calculateFileHash 计算文件哈希值
// @fileBytes 文件字节流
// @return 文件哈希, 错误
func (s *MediaService) calculateFileHash(fileBytes []byte) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, bytes.NewReader(fileBytes)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// generateFileName 生成安全的文件名
// @uploaderID 上传者ID
// @originalName 原始文件名
// @return 生成的文件名
func (s *MediaService) generateFileName(uploaderID uint, originalName string) string {
	ext := strings.ToLower(filepath.Ext(originalName))
	hash := md5.Sum([]byte(fmt.Sprintf("%d-%d-%s", uploaderID, time.Now().UnixNano(), originalName)))
	return hex.EncodeToString(hash[:]) + ext
}

// ensureUploadDir 确保上传目录存在
// @dir 目录路径
// @return 错误
func (s *MediaService) ensureUploadDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建上传目录失败: %v", err)
	}
	return nil
}

// saveFile 保存文件到磁盘
// @filePath 文件路径
// @fileBytes 文件字节流
// @return 错误
func (s *MediaService) saveFile(filePath string, fileBytes []byte) error {
	if err := os.WriteFile(filePath, fileBytes, 0644); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}
	return nil
}

// GetMediaByID 根据ID获取媒体文件
// @mediaID 媒体文件ID
// @return 媒体文件, 错误
func (s *MediaService) GetMediaByID(mediaID uint) (*model.Media, error) {
	media, err := s.mediaRepo.GetMediaByID(mediaID)
	if err != nil {
		return nil, ErrMediaNotFound
	}
	return media, nil
}

// GetMediaList 获取媒体文件列表
// @page 页码
// @limit 每页数量
// @return 媒体文件列表, 总数, 错误
func (s *MediaService) GetMediaList(page, limit int) ([]model.Media, int64, error) {
	return s.mediaRepo.GetMediaList(page, limit)
}

// CheckDeletePermission 检查删除权限
// @mediaID 媒体文件ID
// @userID 用户ID
// @role 用户角色
// @return 错误
func (s *MediaService) CheckDeletePermission(mediaID uint, userID uint, role string) error {
	media, err := s.mediaRepo.GetMediaByID(mediaID)
	if err != nil {
		return ErrMediaNotFound
	}

	// 只有上传者本人或管理员可以删除
	if media.UploaderID != userID && role != "admin" {
		return ErrPermissionDenied
	}

	return nil
}

// DeleteMediaWithFile 删除媒体文件（包括物理文件和数据库记录）
// @mediaID 媒体文件ID
// @return 错误
func (s *MediaService) DeleteMediaWithFile(mediaID uint) error {
	// 获取媒体信息
	media, err := s.mediaRepo.GetMediaByID(mediaID)
	if err != nil {
		return ErrMediaNotFound
	}

	// 先删除数据库记录
	if err := s.mediaRepo.DeleteMedia(mediaID); err != nil {
		return fmt.Errorf("删除数据库记录失败: %v", err)
	}

	// 再删除物理文件（即使失败也不影响数据库操作）
	filePath := filepath.Join(config.Cfg.Media.MediaDir, media.FilePath)
	if err := os.Remove(filePath); err != nil {
		// 记录日志但不返回错误，因为数据库记录已删除
		// TODO: 添加日志记录
		// log.Printf("删除物理文件失败: %v", err)
	}

	return nil
}

// GetMediaFilePath 获取媒体文件的完整路径（用于文件访问）
// @fileName 文件名
// @return 文件路径, 错误
func (s *MediaService) GetMediaFilePath(fileName string) (string, error) {
	// 安全性检查
	safeFileName := filepath.Base(fileName)
	if safeFileName != fileName {
		return "", ErrInvalidFileName
	}

	filePath := filepath.Join(config.Cfg.Media.MediaDir, safeFileName)

	// 确保解析后的路径仍在上传目录内
	absUploadDir, _ := filepath.Abs(config.Cfg.Media.MediaDir)
	absFilePath, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(absFilePath, absUploadDir) {
		return "", ErrPermissionDenied
	}

	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "", ErrMediaNotFound
	}

	// 防止目录遍历
	if fileInfo.IsDir() {
		return "", ErrPermissionDenied
	}

	return filePath, nil
}

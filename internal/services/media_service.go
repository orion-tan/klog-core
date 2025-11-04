package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"klog-backend/internal/config"
	"klog-backend/internal/model"
	"klog-backend/internal/queue"
	"klog-backend/internal/repository"
	"klog-backend/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"gorm.io/gorm"
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
	FileBytes []byte    // 用于小文件或base64上传
	Reader    io.Reader // 用于流式上传（大文件）
	MimeType  string
	Size      int64
}

type MediaService struct {
	mediaRepo *repository.MediaRepository
	fileQueue *queue.FileQueue
}

func NewMediaService(mediaRepo *repository.MediaRepository, fileQueue *queue.FileQueue) *MediaService {
	return &MediaService{
		mediaRepo: mediaRepo,
		fileQueue: fileQueue,
	}
}

// SaveMediaFile 保存媒体文件（包括物理文件和数据库记录）
// @fileData 文件数据
// @return 媒体文件, 错误
func (s *MediaService) SaveMediaFile(fileData *FileData) (*model.Media, error) {
	// 1. 验证文件
	if err := s.validateFile(fileData); err != nil {
		return nil, err
	}

	// 2. 生成唯一文件名
	fileName := s.generateFileName(fileData.FileName)

	// 3. 确保上传目录存在
	uploadDir := config.Cfg.Media.MediaDir
	if err := s.ensureUploadDir(uploadDir); err != nil {
		return nil, err
	}

	// 4. 保存物理文件（流式处理）
	filePath := filepath.Join(uploadDir, fileName)
	fileHash, err := s.saveFileStream(filePath, fileData)
	if err != nil {
		return nil, err
	}

	// 5. 检查文件是否已存在（去重）
	existingMedia, err := s.mediaRepo.GetMediaByHash(fileHash)
	if err == nil && existingMedia != nil {
		// 文件已存在，删除刚保存的文件，返回已存在的记录
		_ = os.Remove(filePath)
		return existingMedia, nil
	}

	// 6. 创建数据库记录
	media := &model.Media{
		FileName: fileData.FileName,
		FilePath: fileName,
		FileHash: fileHash,
		MimeType: fileData.MimeType,
		Size:     fileData.Size,
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

	// 检查客户端提供的 MIME 类型
	if !AllowedMimeTypes[fileData.MimeType] {
		return fmt.Errorf("%w: %s", ErrInvalidFileType, fileData.MimeType)
	}

	// 二次验证：检测文件实际内容的 MIME 类型（防止伪造）
	// 针对流式上传检测
	if fileData.Reader != nil {
		detectedMimeType, fullReader := s.detectMimeTypeByReader(fileData.Reader)
		fileData.Reader = fullReader
		if !AllowedMimeTypes[detectedMimeType] {
			return fmt.Errorf("%w: 检测到的文件类型为 %s，不允许上传", ErrInvalidFileType, detectedMimeType)
		}
	} else {
		detectedMimeType := s.detectMimeType(fileData.FileBytes)
		if !AllowedMimeTypes[detectedMimeType] {
			return fmt.Errorf("%w: 检测到的文件类型为 %s，不允许上传", ErrInvalidFileType, detectedMimeType)
		}
	}

	return nil
}

// detectMimeType 检测文件实际的 MIME 类型
// @fileBytes 文件字节流
// @return MIME 类型
func (s *MediaService) detectMimeType(fileBytes []byte) string {
	mime := mimetype.Detect(fileBytes)
	return mime.String()
}

// detectMimeType 检测文件实际的 MIME 类型
// @fileReader 文件读取器
// @return MIME 类型
func (s *MediaService) detectMimeTypeByReader(fileReader io.Reader) (string, io.Reader) {
	header := make([]byte, 1024)
	n, err := io.ReadFull(fileReader, header)
	// 如果文件本身小于512字节，io.ReadFull会返回ErrUnexpectedEOF
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return "", fileReader
	}
	detectedMimeType := mimetype.Detect(header[:n])
	fullReader := io.MultiReader(bytes.NewReader(header[:n]), fileReader)
	return detectedMimeType.String(), fullReader
}

// generateFileName 生成安全的文件名（使用 SHA-256）
// @originalName 原始文件名
// @return 生成的文件名
func (s *MediaService) generateFileName(originalName string) string {
	ext := strings.ToLower(filepath.Ext(originalName))
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d-%s", time.Now().UnixNano(), originalName)))
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

// saveFileStream 流式保存文件到磁盘（支持大文件）
// @filePath 文件路径
// @fileData 文件数据
// @return 文件哈希, 错误
func (s *MediaService) saveFileStream(filePath string, fileData *FileData) (string, error) {
	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	// 创建哈希计算器
	hash := sha256.New()

	// 使用 MultiWriter 同时写入文件和计算哈希
	writer := io.MultiWriter(dst, hash)

	var written int64

	// 根据数据来源选择处理方式
	if fileData.Reader != nil {
		// 流式写入
		written, err = io.Copy(writer, fileData.Reader)
	} else if fileData.FileBytes != nil {
		// 字节数组上传（小文件或base64）
		written, err = io.Copy(writer, bytes.NewReader(fileData.FileBytes))
	} else {
		return "", errors.New("无效的文件数据：既没有Reader也没有FileBytes")
	}

	if err != nil {
		_ = os.Remove(filePath)
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	// 验证写入的大小
	if written != fileData.Size {
		_ = os.Remove(filePath)
		return "", fmt.Errorf("文件大小不匹配：期望%d字节，实际写入%d字节", fileData.Size, written)
	}

	// 返回文件哈希
	fileHash := hex.EncodeToString(hash.Sum(nil))
	return fileHash, nil
}

// GetMediaByID 根据ID获取媒体文件
// @mediaID 媒体文件ID
// @return 媒体文件, 错误
func (s *MediaService) GetMediaByID(mediaID uint) (*model.Media, error) {
	media, err := s.mediaRepo.GetMediaByID(mediaID)
	if err != nil {
		// 精确判断错误类型
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMediaNotFound
		}
		return nil, fmt.Errorf("查询媒体文件失败: %w", err)
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

// DeleteMediaWithFile 删除媒体文件（包括物理文件和数据库记录）
// @mediaID 媒体文件ID
// @return 错误
func (s *MediaService) DeleteMediaWithFile(mediaID uint) error {
	media, err := s.mediaRepo.GetMediaByID(mediaID)
	if err != nil {
		// 精确判断错误类型
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMediaNotFound
		}
		return fmt.Errorf("查询媒体文件失败: %w", err)
	}

	if err := s.mediaRepo.DeleteMedia(mediaID); err != nil {
		return fmt.Errorf("删除数据库记录失败: %v", err)
	}

	// 异步删除物理文件（提交到消息队列处理）
	filePath := filepath.Join(config.Cfg.Media.MediaDir, media.FilePath)
	if err := s.fileQueue.PublishDeleteTask(filePath); err != nil {
		// 提交失败记录日志，定时清理任务会处理孤儿文件
		utils.SugarLogger.Warnf("提交文件删除任务失败: %s, 错误: %v", filePath, err)
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

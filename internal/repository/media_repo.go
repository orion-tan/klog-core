package repository

import (
	"klog-backend/internal/model"

	"gorm.io/gorm"
)

type MediaRepository struct {
	DB *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{DB: db}
}

// 开启事务
func (r *MediaRepository) WithTransaction(fn func(*gorm.DB) error) error {
	return r.DB.Transaction(fn)
}

// CreateMedia 创建媒体文件记录
// @media 媒体文件
// @return 错误
func (r *MediaRepository) CreateMedia(media *model.Media) error {
	return r.DB.Create(media).Error
}

// GetMediaByID 根据ID获取媒体文件
// @mediaID 媒体文件ID
// @return 媒体文件, 错误
func (r *MediaRepository) GetMediaByID(mediaID uint) (*model.Media, error) {
	var media model.Media
	err := r.DB.First(&media, mediaID).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

// GetMediaList 获取媒体文件列表（带分页）
// @page 页码
// @limit 每页数量
// @return 媒体文件列表, 总数, 错误
func (r *MediaRepository) GetMediaList(page, limit int) ([]model.Media, int64, error) {
	var mediaList []model.Media
	var total int64

	query := r.DB.Model(&model.Media{})

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit).Order("created_at DESC")

	if err := query.Find(&mediaList).Error; err != nil {
		return nil, 0, err
	}

	return mediaList, total, nil
}

// GetMediaByHash 根据文件哈希获取媒体文件（用于去重）
// @fileHash 文件哈希值
// @return 媒体文件, 错误
func (r *MediaRepository) GetMediaByHash(fileHash string) (*model.Media, error) {
	var media model.Media
	err := r.DB.Where("file_hash = ?", fileHash).First(&media).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

// DeleteMedia 删除媒体文件
// @mediaID 媒体文件ID
// @return 错误
func (r *MediaRepository) DeleteMedia(mediaID uint) error {
	return r.DB.Delete(&model.Media{}, mediaID).Error
}

// 在事务中删除媒体文件记录
// @tx 事务
// @mediaID 媒体文件ID
// @return 错误
func (r *MediaRepository) DeleteMediaInTx(tx *gorm.DB, mediaID uint) error {
	return tx.Delete(&model.Media{}, mediaID).Error
}

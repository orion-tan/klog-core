package services

import (
	"errors"
	"fmt"
	"klog-backend/internal/api"
	"klog-backend/internal/cache"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"
	"time"

	"gorm.io/gorm"
)

type TagService struct {
	tagRepo  *repository.TagRepository
	postRepo *repository.PostRepository
}

func NewTagService(tagRepo *repository.TagRepository, postRepo *repository.PostRepository) *TagService {
	return &TagService{
		tagRepo:  tagRepo,
		postRepo: postRepo,
	}
}

const (
	tagCacheKey = "tags:all"
	tagCacheTTL = 10 * time.Minute
)

// GetTags 获取所有标签（带缓存）
func (s *TagService) GetTags() ([]model.Tag, error) {
	// 尝试从缓存获取
	var tags []model.Tag
	err := cache.Get(tagCacheKey, &tags)
	if err == nil {
		return tags, nil
	}

	// 缓存未命中或Redis未启用，从数据库获取
	tags, err = s.tagRepo.GetTags()
	if err != nil {
		return nil, err
	}

	// 写入缓存
	_ = cache.Set(tagCacheKey, tags, tagCacheTTL)

	return tags, nil
}

// CreateTag 创建标签
func (s *TagService) CreateTag(req *api.TagCreateRequest) (*model.Tag, error) {
	// 检查slug是否已存在
	_, err := s.tagRepo.GetTagBySlug(req.Slug)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("标签slug已存在")
	}

	tag := &model.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := s.tagRepo.CreateTag(tag); err != nil {
		return nil, err
	}

	// 清除缓存
	_ = cache.Delete(tagCacheKey)

	return tag, nil
}

// UpdateTag 更新标签
func (s *TagService) UpdateTag(id uint, req *api.TagUpdateRequest) (*model.Tag, error) {
	tag, err := s.tagRepo.GetTagByID(id)
	if err != nil {
		return nil, errors.New("标签不存在")
	}

	// 检查slug是否被其他标签使用
	if req.Slug != "" && req.Slug != tag.Slug {
		existingTag, err := s.tagRepo.GetTagBySlug(req.Slug)
		if err == nil && existingTag.ID != id {
			return nil, errors.New("标签slug已被使用")
		}
	}

	// 更新字段
	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Slug != "" {
		tag.Slug = req.Slug
	}

	if err := s.tagRepo.UpdateTag(tag); err != nil {
		return nil, err
	}

	// 清除缓存
	_ = cache.Delete(tagCacheKey)

	return tag, nil
}

// DeleteTag 删除标签
func (s *TagService) DeleteTag(id uint) error {
	_, err := s.tagRepo.GetTagByID(id)
	if err != nil {
		return errors.New("标签不存在")
	}

	// 检查是否有文章使用该标签
	count, err := s.postRepo.CountPostsByTagID(id)
	if err != nil {
		return fmt.Errorf("检查标签使用情况失败: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("该标签被 %d 篇文章使用，无法删除", count)
	}

	if err := s.tagRepo.DeleteTag(id); err != nil {
		return err
	}

	// 清除缓存
	_ = cache.Delete(tagCacheKey)

	return nil
}

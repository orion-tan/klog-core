package services

import (
	"errors"
	"klog-backend/internal/api"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"

	"gorm.io/gorm"
)

type TagService struct {
	tagRepo  *repository.TagRepository
	postRepo *repository.PostRepository
}

func NewTagService(tagRepo *repository.TagRepository, postRepo *repository.PostRepository) *TagService {
	return &TagService{tagRepo: tagRepo, postRepo: postRepo}
}

// CreateTag 创建标签
// @req 创建标签请求
// @return 标签, 错误
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

	return tag, nil
}

// GetTagByID 根据ID获取标签
// @tagID 标签ID
// @return 标签, 错误
func (s *TagService) GetTagByID(tagID uint) (*model.Tag, error) {
	return s.tagRepo.GetTagByID(tagID)
}

// GetTags 获取所有标签
// @return 标签列表, 错误
func (s *TagService) GetTags() ([]model.Tag, error) {
	return s.tagRepo.GetTags()
}

// UpdateTag 更新标签
// @tagID 标签ID
// @req 更新标签请求
// @return 标签, 错误
func (s *TagService) UpdateTag(tagID uint, req *api.TagUpdateRequest) (*model.Tag, error) {
	tag, err := s.tagRepo.GetTagByID(tagID)
	if err != nil {
		return nil, errors.New("标签不存在")
	}

	// 检查slug是否已被其他标签使用
	if req.Slug != "" && req.Slug != tag.Slug {
		existingTag, err := s.tagRepo.GetTagBySlug(req.Slug)
		if err == nil && existingTag.ID != tagID {
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

	return tag, nil
}

// DeleteTag 删除标签
// @tagID 标签ID
// @return 错误
func (s *TagService) DeleteTag(tagID uint) error {
	_, err := s.tagRepo.GetTagByID(tagID)
	if err != nil {
		return errors.New("标签不存在")
	}

	count, err := s.postRepo.CountPostsByTagID(tagID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("标签下有文章，不能删除")
	}
	return s.tagRepo.DeleteTag(tagID)
}

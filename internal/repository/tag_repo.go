package repository

import (
	"klog-backend/internal/api"
	"klog-backend/internal/model"

	"gorm.io/gorm"
)

type TagRepository struct {
	DB *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{DB: db}
}

// CreateTag 创建标签
// @tag 标签
// @return 错误
func (r *TagRepository) CreateTag(tag *model.Tag) error {
	return r.DB.Create(tag).Error
}

// GetTagByID 根据ID获取标签
// @tagID 标签ID
// @return 标签, 错误
func (r *TagRepository) GetTagByID(tagID uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.DB.First(&tag, tagID).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetTagBySlug 根据Slug获取标签
// @slug 标签Slug
// @return 标签, 错误
func (r *TagRepository) GetTagBySlug(slug string) (*model.Tag, error) {
	var tag model.Tag
	err := r.DB.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetTagByName 根据Name获取标签
// @name 标签Name
// @return 标签, 错误
func (r *TagRepository) GetTagByName(name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.DB.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetTags 获取所有标签
// @return 标签列表, 错误
func (r *TagRepository) GetTags() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.DB.Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// GetTagsWithCount 获取所有标签及其文章数量
// @return 标签列表（含文章数量）, 错误
func (r *TagRepository) GetTagsWithCount() ([]api.TagWithCount, error) {
	var tags []api.TagWithCount
	err := r.DB.Table("tags").
		Select("tags.id, tags.name, tags.slug, COUNT(DISTINCT post_tags.post_id) as post_count").
		Joins("LEFT JOIN post_tags ON tags.id = post_tags.tag_id").
		Joins("LEFT JOIN posts ON post_tags.post_id = posts.id AND posts.status = ?", "published").
		Group("tags.id").
		Order("tags.id ASC").
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// GetTagsBySlugs 根据Slug列表获取标签列表
// @slugs Slug列表
// @return 标签列表, 错误
func (r *TagRepository) GetTagsBySlugs(slugs []string) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.DB.Where("slug IN ?", slugs).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// UpdateTag 更新标签
// @tag 标签
// @return 错误
func (r *TagRepository) UpdateTag(tag *model.Tag) error {
	return r.DB.Save(tag).Error
}

// DeleteTag 删除标签
// @tagID 标签ID
// @return 错误
func (r *TagRepository) DeleteTag(tagID uint) error {
	return r.DB.Delete(&model.Tag{}, tagID).Error
}

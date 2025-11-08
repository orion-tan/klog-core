package repository

import (
	"klog-backend/internal/api"
	"klog-backend/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

// CreateCategory 创建分类
// @category 分类
// @return 错误
func (r *CategoryRepository) CreateCategory(category *model.Category) error {
	return r.DB.Create(category).Error
}

// GetCategoryByID 根据ID获取分类
// @categoryID 分类ID
// @return 分类, 错误
func (r *CategoryRepository) GetCategoryByID(categoryID uint) (*model.Category, error) {
	var category model.Category
	err := r.DB.First(&category, categoryID).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetCategoryBySlug 根据Slug获取分类
// @slug 分类Slug
// @return 分类, 错误
func (r *CategoryRepository) GetCategoryBySlug(slug string) (*model.Category, error) {
	var category model.Category
	err := r.DB.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetCategories 获取所有分类
// @return 分类列表, 错误
func (r *CategoryRepository) GetCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.DB.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoriesWithCount 获取所有分类及其文章数量
// @return 分类列表（含文章数量）, 错误
func (r *CategoryRepository) GetCategoriesWithCount() ([]api.CategoryWithCount, error) {
	var categories []api.CategoryWithCount
	err := r.DB.Table("categories").
		Select("categories.id, categories.name, categories.slug, categories.description, COUNT(posts.id) as post_count").
		Joins("LEFT JOIN posts ON categories.id = posts.category_id AND posts.status = ?", "published").
		Group("categories.id").
		Order("categories.id ASC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// UpdateCategory 更新分类
// @category 分类
// @return 错误
func (r *CategoryRepository) UpdateCategory(category *model.Category) error {
	return r.DB.Save(category).Error
}

// DeleteCategory 删除分类
// @categoryID 分类ID
// @return 错误
func (r *CategoryRepository) DeleteCategory(categoryID uint) error {
	return r.DB.Delete(&model.Category{}, categoryID).Error
}


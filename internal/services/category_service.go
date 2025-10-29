package services

import (
	"errors"
	"klog-backend/internal/api"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"

	"gorm.io/gorm"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	postRepo     *repository.PostRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository, postRepo *repository.PostRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo, postRepo: postRepo}
}

// CreateCategory 创建分类
// @req 创建分类请求
// @return 分类, 错误
func (s *CategoryService) CreateCategory(req *api.CategoryCreateRequest) (*model.Category, error) {
	// 检查slug是否已存在
	_, err := s.categoryRepo.GetCategoryBySlug(req.Slug)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("分类slug已存在")
	}

	category := &model.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := s.categoryRepo.CreateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID 根据ID获取分类
// @categoryID 分类ID
// @return 分类, 错误
func (s *CategoryService) GetCategoryByID(categoryID uint) (*model.Category, error) {
	return s.categoryRepo.GetCategoryByID(categoryID)
}

// GetCategories 获取所有分类
// @return 分类列表, 错误
func (s *CategoryService) GetCategories() ([]model.Category, error) {
	return s.categoryRepo.GetCategories()
}

// UpdateCategory 更新分类
// @categoryID 分类ID
// @req 更新分类请求
// @return 分类, 错误
func (s *CategoryService) UpdateCategory(categoryID uint, req *api.CategoryUpdateRequest) (*model.Category, error) {
	category, err := s.categoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		return nil, errors.New("分类不存在")
	}

	// 检查slug是否已被其他分类使用
	if req.Slug != "" && req.Slug != category.Slug {
		existingCategory, err := s.categoryRepo.GetCategoryBySlug(req.Slug)
		if err == nil && existingCategory.ID != categoryID {
			return nil, errors.New("分类slug已被使用")
		}
	}

	// 更新字段
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := s.categoryRepo.UpdateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory 删除分类
// @categoryID 分类ID
// @return 错误
func (s *CategoryService) DeleteCategory(categoryID uint) error {
	_, err := s.categoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		return errors.New("分类不存在")
	}

	count, err := s.postRepo.CountPostsByCategoryID(categoryID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("分类下有文章，不能删除")
	}
	return s.categoryRepo.DeleteCategory(categoryID)
}

package services

import (
	"errors"
	"fmt"
	"klog-backend/internal/api"
	"klog-backend/internal/cache"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"
	"time"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	postRepo     *repository.PostRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository, postRepo *repository.PostRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		postRepo:     postRepo,
	}
}

const (
	categoryCacheKey = "categories:all"
	cacheTTL         = 10 * time.Minute
)

// GetCategories 获取所有分类（带缓存）
func (s *CategoryService) GetCategories() ([]model.Category, error) {
	// 尝试从缓存获取
	var categories []model.Category
	err := cache.Get(categoryCacheKey, &categories)
	if err == nil {
		return categories, nil
	}

	// 缓存未命中或Redis未启用，从数据库获取
	categories, err = s.categoryRepo.GetCategories()
	if err != nil {
		return nil, err
	}

	// 写入缓存
	_ = cache.Set(categoryCacheKey, categories, cacheTTL)

	return categories, nil
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(req *api.CategoryCreateRequest) (*model.Category, error) {
	// 检查slug是否已存在
	category, err := s.categoryRepo.GetCategoryBySlug(req.Slug)
	if err == nil && category.ID != 0 {
		return nil, errors.New("分类slug已存在")
	}

	category = &model.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: &req.Description,
	}

	if err := s.categoryRepo.CreateCategory(category); err != nil {
		return nil, err
	}

	// 清除缓存
	_ = cache.Delete(categoryCacheKey)

	return category, nil
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(id uint, req *api.CategoryUpdateRequest) (*model.Category, error) {
	category, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return nil, errors.New("分类不存在")
	}

	// 检查slug是否被其他分类使用
	if req.Slug != nil && *req.Slug != "" && *req.Slug != category.Slug {
		existingCategory, err := s.categoryRepo.GetCategoryBySlug(*req.Slug)
		if err == nil && existingCategory.ID != id {
			return nil, errors.New("分类slug已被使用")
		}
	}

	// 更新字段
	if req.Name != nil && *req.Name != "" {
		category.Name = *req.Name
	}
	if req.Slug != nil && *req.Slug != "" {
		category.Slug = *req.Slug
	}
	if req.Description != nil {
		category.Description = req.Description
	}

	if err := s.categoryRepo.UpdateCategory(category); err != nil {
		return nil, err
	}

	// 清除缓存
	_ = cache.Delete(categoryCacheKey)

	return category, nil
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(id uint) error {
	_, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return errors.New("分类不存在")
	}

	// 检查是否有文章使用该分类
	count, err := s.postRepo.CountPostsByCategoryID(id)
	if err != nil {
		return fmt.Errorf("查询数据库失败: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("该分类下有 %d 篇文章，无法删除", count)
	}

	if err := s.categoryRepo.DeleteCategory(id); err != nil {
		return err
	}

	// 清除缓存
	_ = cache.Delete(categoryCacheKey)

	return nil
}

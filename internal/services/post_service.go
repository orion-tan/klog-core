package services

import (
	"errors"
	"fmt"
	"klog-backend/internal/api"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"
	"klog-backend/internal/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PostService struct {
	postRepo     *repository.PostRepository
	tagRepo      *repository.TagRepository
	categoryRepo *repository.CategoryRepository
}

func NewPostService(postRepo *repository.PostRepository, tagRepo *repository.TagRepository, categoryRepo *repository.CategoryRepository) *PostService {
	return &PostService{
		postRepo:     postRepo,
		tagRepo:      tagRepo,
		categoryRepo: categoryRepo,
	}
}

// CreatePost 创建文章
// @req 创建文章请求
// @authorID 作者ID
// @return 文章, 错误
func (s *PostService) CreatePost(req *api.PostCreateRequest, authorID uint) (*model.Post, error) {
	// 检查slug是否已存在
	_, err := s.postRepo.GetPostBySlug(req.Slug, false)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("文章slug已存在")
	}

	// 检查标签
	var tags []model.Tag
	if len(req.Tags) > 0 {
		tags, err = s.tagRepo.GetTagsBySlugs(req.Tags)
		if err != nil {
			return nil, errors.New("获取标签失败")
		}

		if len(tags) != len(req.Tags) {
			foundSlugs := make(map[string]bool)
			for _, tag := range tags {
				foundSlugs[tag.Slug] = true
			}
			missingSlugs := []string{}
			for _, slug := range req.Tags {
				if !foundSlugs[slug] {
					missingSlugs = append(missingSlugs, slug)
				}
			}
			return nil, fmt.Errorf("以下标签不存在: %s", strings.Join(missingSlugs, ", "))
		}
	}

	// 检查分类是否存在
	if req.CategoryID != nil {
		_, err := s.categoryRepo.GetCategoryByID(*req.CategoryID)
		if err != nil {
			return nil, errors.New("分类不存在")
		}
	}

	var post *model.Post

	err = s.postRepo.WithTransaction(func(tx *gorm.DB) error {

		// 构建文章对象
		post = &model.Post{
			CategoryID:    req.CategoryID,
			AuthorID:      authorID,
			Title:         req.Title,
			Slug:          req.Slug,
			Content:       req.Content,
			Excerpt:       req.Excerpt,
			CoverImageURL: req.CoverImageURL,
			Status:        req.Status,
		}

		// 如果状态为published，设置发布时间
		if req.Status == "published" {
			now := time.Now()
			post.PublishedAt = &now
		}

		if err := s.postRepo.CreatePostInTx(tx, post); err != nil {
			return err
		}

		// 关联标签
		if len(tags) > 0 {
			if err := s.postRepo.AssociateTagsInTx(tx, post, tags); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		utils.SugarLogger.Errorln("创建文章失败，自动回滚：", err)
		return nil, err
	}

	return s.postRepo.GetPostByID(post.ID, true)
}

// GetPostByID 根据ID获取文章
// @postID 文章ID
// @return 文章, 错误
func (s *PostService) GetPostByID(postID uint) (*model.Post, error) {
	post, err := s.postRepo.GetPostByID(postID, true)
	if err != nil {
		return nil, err
	}

	// 增加浏览量
	// _ = s.postRepo.IncrementViewCount(postID)

	return post, nil
}

// GetPosts 获取文章列表
// @page 页码
// @limit 每页数量
// @status 文章状态
// @authorID 作者ID
// @categorySlug 分类Slug
// @tagSlug 标签Slug
// @sortBy 排序字段
// @order 排序方式
// @return 文章列表, 总数, 错误
func (s *PostService) GetPosts(page, limit int, status, categorySlug, tagSlug, sortBy, order string, authorID *uint) ([]model.Post, int64, error) {
	return s.postRepo.GetPosts(page, limit, status, categorySlug, tagSlug, sortBy, order, authorID)
}

// UpdatePost 更新文章
// @postID 文章ID
// @req 更新文章请求
// @return 文章, 错误
func (s *PostService) UpdatePost(postID uint, req *api.PostUpdateRequest) (*model.Post, error) {
	post, err := s.postRepo.GetPostByID(postID, false)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 检查slug是否已被其他文章使用
	if req.Slug != "" && req.Slug != post.Slug {
		existingPost, err := s.postRepo.GetPostBySlug(req.Slug, false)
		if err == nil && existingPost.ID != postID {
			return nil, errors.New("文章slug已被使用")
		}
	}

	// 检查分类是否存在
	if req.CategoryID != nil {
		_, err := s.categoryRepo.GetCategoryByID(*req.CategoryID)
		if err != nil {
			return nil, errors.New("分类不存在，请先创建分类")
		}
	}

	// 更新字段
	if req.CategoryID != nil {
		post.CategoryID = req.CategoryID
	}
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Slug != "" {
		post.Slug = req.Slug
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Excerpt != "" {
		post.Excerpt = req.Excerpt
	}
	if req.CoverImageURL != "" {
		post.CoverImageURL = req.CoverImageURL
	}
	if req.Status != "" {
		// 如果从非发布状态更改为发布状态，设置发布时间
		if post.Status != "published" && req.Status == "published" {
			now := time.Now()
			post.PublishedAt = &now
		}
		post.Status = req.Status
	}

	// 检查标签
	var tags []model.Tag
	if len(req.Tags) > 0 {
		tags, err = s.tagRepo.GetTagsBySlugs(req.Tags)
		if err != nil {
			return nil, errors.New("获取标签失败")
		}
		if len(tags) != len(req.Tags) {
			foundSlugs := make(map[string]bool)
			for _, tag := range tags {
				foundSlugs[tag.Slug] = true
			}
			missingSlugs := []string{}
			for _, slug := range req.Tags {
				if !foundSlugs[slug] {
					missingSlugs = append(missingSlugs, slug)
				}
			}
			return nil, fmt.Errorf("以下标签不存在: %s", strings.Join(missingSlugs, ", "))
		}
	}

	// 在事务中更新文章
	err = s.postRepo.WithTransaction(func(tx *gorm.DB) error {
		if err := s.postRepo.UpdatePostInTx(tx, post); err != nil {
			return err
		}

		if err := s.postRepo.AssociateTagsInTx(tx, post, tags); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		utils.SugarLogger.Errorln("更新文章失败，自动回滚：", err)
		return nil, err
	}

	return s.postRepo.GetPostByID(post.ID, true)
}

// DeletePost 删除文章
// @postID 文章ID
// @return 错误
func (s *PostService) DeletePost(postID uint) error {
	_, err := s.postRepo.GetPostByID(postID, false)
	if err != nil {
		return errors.New("文章不存在")
	}
	return s.postRepo.DeletePost(postID)
}

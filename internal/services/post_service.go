package services

import (
	"errors"
	"fmt"
	"klog-backend/internal/api"
	"klog-backend/internal/cache"
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

// generatePostListCacheKey 生成文章列表缓存key
func generatePostListCacheKey(page, limit int, status, categorySlug, tagSlug, sortBy, order string) string {
	return fmt.Sprintf("posts:list:%s:%s:%s:%s:%s:page:%d:limit:%d",
		status, categorySlug, tagSlug, sortBy, order, page, limit)
}

// clearPostListCache 清理所有文章列表缓存
func clearPostListCache() {
	_ = cache.DeleteByPattern("posts:list:*")
	utils.SugarLogger.Info("已清理所有文章列表缓存")
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

	// 清理文章列表缓存
	clearPostListCache()

	return s.postRepo.GetPostByID(post.ID, true)
}

// GetPostByID 根据ID获取文章（带缓存）
// @postID 文章ID
// @return 文章, 错误
func (s *PostService) GetPostByID(postID uint) (*model.Post, error) {
	cacheKey := fmt.Sprintf("post:detail:%d", postID)
	var post model.Post

	// 使用缓存或从数据库加载
	err := cache.GetOrSet(cacheKey, &post, 10*time.Minute, func() (interface{}, error) {
		p, err := s.postRepo.GetPostByID(postID, true)
		if err != nil {
			return nil, err
		}
		return p, nil
	})

	if err != nil {
		return nil, err
	}

	return &post, nil
}

// GetPosts 获取文章列表（带缓存）
// @page 页码
// @limit 每页数量
// @status 文章状态
// @authorID 作者ID
// @categorySlug 分类Slug
// @tagSlug 标签Slug
// @sortBy 排序字段
// @order 排序方式
// @return 文章列表, 总数, 错误
func (s *PostService) GetPosts(page, limit int, status, categorySlug, tagSlug, sortBy, order string) ([]model.Post, int64, error) {
	// 生成缓存key
	cacheKey := generatePostListCacheKey(page, limit, status, categorySlug, tagSlug, sortBy, order)

	// 定义缓存数据结构
	type CachedPostList struct {
		Posts []model.Post `json:"posts"`
		Total int64        `json:"total"`
	}

	var cached CachedPostList

	// 尝试从缓存获取或从数据库加载
	err := cache.GetOrSet(cacheKey, &cached, 5*time.Minute, func() (interface{}, error) {
		posts, total, err := s.postRepo.GetPosts(page, limit, status, categorySlug, tagSlug, sortBy, order)
		if err != nil {
			return nil, err
		}
		return &CachedPostList{
			Posts: posts,
			Total: total,
		}, nil
	})

	if err != nil {
		return nil, 0, err
	}

	return cached.Posts, cached.Total, nil
}

// GetPostsByCursor 使用游标分页获取文章列表
// @cursorStr 游标字符串（空表示首页）
// @limit 每页数量
// @status 文章状态
// @categorySlug 分类Slug
// @tagSlug 标签Slug
// @sortBy 排序字段
// @order 排序方式
// @return 文章列表, 下一页游标, 是否有更多数据, 错误
func (s *PostService) GetPostsByCursor(cursorStr string, limit int, status, categorySlug, tagSlug, sortBy, order string) ([]model.Post, *string, bool, error) {
	// 解码游标
	var cursor *utils.CursorData
	var err error
	if cursorStr != "" {
		cursor, err = utils.DecodeCursor(cursorStr)
		if err != nil {
			return nil, nil, false, err
		}
	}

	// 从数据库获取数据
	posts, hasMore, err := s.postRepo.GetPostsByCursor(cursor, limit, status, categorySlug, tagSlug, sortBy, order)
	if err != nil {
		return nil, nil, false, err
	}

	// 生成下一页游标
	var nextCursor *string
	if hasMore && len(posts) > 0 {
		lastPost := posts[len(posts)-1]

		// 根据排序字段获取值
		var sortValue string
		switch sortBy {
		case "published_at":
			sortValue = utils.FormatSortValue(sortBy, lastPost.PublishedAt)
		case "created_at":
			sortValue = utils.FormatSortValue(sortBy, lastPost.CreatedAt)
		case "updated_at":
			sortValue = utils.FormatSortValue(sortBy, lastPost.UpdatedAt)
		case "view_count":
			sortValue = utils.FormatSortValue(sortBy, lastPost.ViewCount)
		case "title":
			sortValue = utils.FormatSortValue(sortBy, lastPost.Title)
		default:
			sortValue = utils.FormatSortValue("published_at", lastPost.PublishedAt)
		}

		cursorData := utils.CursorData{
			SortField: sortBy,
			SortValue: sortValue,
			ID:        lastPost.ID,
		}
		encoded := utils.EncodeCursor(cursorData)
		nextCursor = &encoded
	}

	return posts, nextCursor, hasMore, nil
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
	if req.Slug != nil && *req.Slug != "" && *req.Slug != post.Slug {
		existingPost, err := s.postRepo.GetPostBySlug(*req.Slug, false)
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
		post.CategoryID = req.CategoryID
	}

	// 更新字段
	if req.Title != nil && *req.Title != "" {
		post.Title = *req.Title
	}
	if req.Slug != nil && *req.Slug != "" {
		post.Slug = *req.Slug
	}
	if req.Content != nil && *req.Content != "" {
		post.Content = *req.Content
	}
	if req.Excerpt != nil {
		post.Excerpt = *req.Excerpt
	}
	if req.CoverImageURL != nil {
		if *req.CoverImageURL != "" {
			if !utils.ValidateURL(*req.CoverImageURL) {
				return nil, errors.New("封面图片URL格式不正确")
			}
			post.CoverImageURL = *req.CoverImageURL
		} else {
			post.CoverImageURL = ""
		}
	}
	if req.Status != nil && *req.Status != "" {
		// 如果从非发布状态更改为发布状态，设置发布时间
		if post.Status != "published" && *req.Status == "published" {
			now := time.Now()
			post.PublishedAt = &now
		}
		// 安全设置状态
		if *req.Status != "draft" && *req.Status != "published" && *req.Status != "archived" {
			return nil, errors.New("文章状态只能是draft、published或archived")
		}
		post.Status = *req.Status
	}

	// 检查标签
	var tags []model.Tag
	if req.Tags != nil && len(*req.Tags) > 0 {
		tags, err = s.tagRepo.GetTagsBySlugs(*req.Tags)
		if err != nil {
			return nil, errors.New("获取标签失败")
		}
		if len(tags) != len(*req.Tags) {
			foundSlugs := make(map[string]bool)
			for _, tag := range tags {
				foundSlugs[tag.Slug] = true
			}
			missingSlugs := []string{}
			for _, slug := range *req.Tags {
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

	// 清理缓存
	clearPostListCache()
	_ = cache.Delete(fmt.Sprintf("post:detail:%d", post.ID))

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

	err = s.postRepo.DeletePost(postID)
	if err != nil {
		return err
	}

	// 清理缓存
	clearPostListCache()
	_ = cache.Delete(fmt.Sprintf("post:detail:%d", postID))

	return nil
}

package repository

import (
	"klog-backend/internal/model"

	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) WithTransaction(fn func(*gorm.DB) error) error {
	return r.DB.Transaction(fn)
}

// CreatePost 在事务中创建文章
// @post 文章
// @return 错误
func (r *PostRepository) CreatePostInTx(tx *gorm.DB, post *model.Post) error {
	return tx.Create(post).Error
}

// CreatePost 创建文章
// @post 文章
// @return 错误
func (r *PostRepository) CreatePost(post *model.Post) error {
	return r.DB.Create(post).Error
}

// 获取某分类下文章数量
func (r *PostRepository) CountPostsByCategoryID(categoryID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Post{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

// 获取某标签下文章数量
func (r *PostRepository) CountPostsByTagID(tagID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Post{}).
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id = ?", tagID).
		Count(&count).Error
	return count, err
}

// GetPostByID 根据ID获取文章
// @postID 文章ID
// @preload 是否预加载关联
// @return 文章, 错误
func (r *PostRepository) GetPostByID(postID uint, preload bool) (*model.Post, error) {
	var post model.Post
	query := r.DB
	if preload {
		query = query.Preload("Category").Preload("Author").Preload("Tags").Preload("Comments")
	}
	err := query.First(&post, postID).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostBySlug 根据Slug获取文章
// @slug 文章Slug
// @preload 是否预加载关联
// @return 文章, 错误
func (r *PostRepository) GetPostBySlug(slug string, preload bool) (*model.Post, error) {
	var post model.Post
	query := r.DB
	if preload {
		query = query.Preload("Category").Preload("Author").Preload("Tags")
	}
	err := query.Where("slug = ?", slug).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPosts 获取文章列表（带分页和筛选）
// @page 页码
// @limit 每页数量
// @status 文章状态
// @authorID 作者ID
// @categorySlug 分类Slug
// @tagSlug 标签Slug
// @sortBy 排序字段
// @order 排序方式
// @return 文章列表, 总数, 错误
func (r *PostRepository) GetPosts(page, limit int, status, categorySlug, tagSlug, sortBy, order string) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.DB.Model(&model.Post{})

	// 筛选条件
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if categorySlug != "" {
		query = query.Joins("JOIN categories ON categories.id = posts.category_id").
			Where("categories.slug = ?", categorySlug)
	}
	if tagSlug != "" {
		query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.slug = ?", tagSlug).
			Distinct()
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序 - 使用白名单防止 SQL 注入
	allowedSortFields := map[string]string{
		"published_at": "published_at",
		"created_at":   "created_at",
		"updated_at":   "updated_at",
		"view_count":   "view_count",
		"title":        "title",
	}
	allowedOrders := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}

	// 验证并设置排序字段
	sortField, ok := allowedSortFields[sortBy]
	if !ok {
		sortField = "published_at"
	}

	// 验证并设置排序方向
	sortOrder, ok := allowedOrders[order]
	if !ok {
		sortOrder = "DESC"
	}

	query = query.Order(sortField + " " + sortOrder)

	// 分页
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// 预加载关联
	query = query.Preload("Category").Preload("Author").Preload("Tags")

	if err := query.Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// UpdatePost 更新文章
// @post 文章
// @return 错误
func (r *PostRepository) UpdatePost(post *model.Post) error {
	return r.DB.Save(post).Error
}

// UpdatePostInTx 在事务中更新文章
// @tx 事务
// @post 文章
// @return 错误
func (r *PostRepository) UpdatePostInTx(tx *gorm.DB, post *model.Post) error {
	return tx.Save(post).Error
}

// DeletePost 删除文章
// @postID 文章ID
// @return 错误
func (r *PostRepository) DeletePost(postID uint) error {
	return r.DB.Select("Comments", "Tags").Delete(&model.Post{}, postID).Error
}

// IncrementViewCount 增加文章浏览量
// @postID 文章ID
// @return 错误
func (r *PostRepository) IncrementViewCount(postID uint, increment uint64) error {
	return r.DB.Model(&model.Post{}).Where("id = ?", postID).Update("view_count", gorm.Expr("view_count + ?", increment)).Error
}

// AssociateTags 关联标签
// @post 文章
// @tags 标签列表
// @return 错误
func (r *PostRepository) AssociateTags(post *model.Post, tags []model.Tag) error {
	return r.DB.Model(post).Association("Tags").Replace(tags)
}

// AssociateTagsInTx 在事务中关联标签
// @post 文章
// @tags 标签列表
// @return 错误
func (r *PostRepository) AssociateTagsInTx(tx *gorm.DB, post *model.Post, tags []model.Tag) error {
	return tx.Model(post).Association("Tags").Replace(tags)
}

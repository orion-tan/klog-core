package repository

import (
	"klog-backend/internal/model"

	"gorm.io/gorm"
)

type CommentRepository struct {
	DB *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

// CreateComment 创建评论
// @comment 评论
// @return 错误
func (r *CommentRepository) CreateComment(comment *model.Comment) error {
	return r.DB.Create(comment).Error
}

// GetCommentByID 根据ID获取评论
// @commentID 评论ID
// @return 评论, 错误
func (r *CommentRepository) GetCommentByID(commentID uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.DB.Preload("User").First(&comment, commentID).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetCommentsByPostID 根据文章ID获取评论列表（树形结构）
// @postID 文章ID
// @return 评论列表, 错误
func (r *CommentRepository) GetCommentsByPostID(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	// 获取文章下所有评论，然后组装树形结构，提高数据库查询性能
	err := r.DB.Where("post_id = ?", postID).
		Preload("User").
		Preload("Replies", "status = ?", "approved").
		Preload("Replies.User").
		Order("created_at DESC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	// 组装树形结构
	tree := make(map[uint][]model.Comment)
	for _, comment := range comments {
		if comment.ParentID == nil {
			tree[comment.ID] = append(tree[comment.ID], comment)
		} else {
			tree[*comment.ParentID] = append(tree[*comment.ParentID], comment)
		}
	}
	return comments, nil
}

// UpdateComment 更新评论
// @comment 评论
// @return 错误
func (r *CommentRepository) UpdateComment(comment *model.Comment) error {
	return r.DB.Save(comment).Error
}

// DeleteComment 删除评论
// @commentID 评论ID
// @return 错误
func (r *CommentRepository) DeleteComment(commentID uint) error {
	return r.DB.Delete(&model.Comment{}, commentID).Error
}

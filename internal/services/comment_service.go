package services

import (
	"errors"
	"klog-backend/internal/api"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository
}

func NewCommentService(commentRepo *repository.CommentRepository, postRepo *repository.PostRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

// CreateComment 创建评论
// @postID 文章ID
// @req 创建评论请求
// @userID 用户ID（可为nil，表示游客）
// @ip IP地址
// @return 评论, 错误
func (s *CommentService) CreateComment(postID uint, req *api.CommentCreateRequest, userID *uint, ip string) (*model.Comment, error) {
	// 检查文章是否存在
	_, err := s.postRepo.GetPostByID(postID, false)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 验证游客评论必填字段
	if userID == nil {
		if req.Name == "" || req.Email == "" {
			return nil, errors.New("游客评论必须填写姓名和邮箱")
		}
	}

	// 验证父评论是否存在且属于同一篇文章
	if req.ParentID != nil {
		parent, err := s.commentRepo.GetCommentByID(*req.ParentID)
		if err != nil {
			return nil, errors.New("父评论不存在")
		}
		if parent.PostID != postID {
			return nil, errors.New("父评论不属于同一篇文章")
		}
	}

	comment := &model.Comment{
		PostID:   postID,
		UserID:   userID,
		Name:     req.Name,
		Email:    req.Email,
		Content:  req.Content,
		IP:       ip,
		ParentID: req.ParentID,
		Status:   "pending", // 默认待审核
	}

	// 如果是已登录用户，自动批准
	if userID != nil {
		comment.Status = "approved"
	}

	if err := s.commentRepo.CreateComment(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

// GetCommentsByPostID 根据文章ID获取评论列表
// @postID 文章ID
// @return 评论列表, 错误
func (s *CommentService) GetCommentsByPostID(postID uint) ([]model.Comment, error) {
	return s.commentRepo.GetCommentsByPostID(postID)
}

// UpdateCommentStatus 更新评论状态
// @commentID 评论ID
// @req 更新评论请求
// @return 评论, 错误
func (s *CommentService) UpdateCommentStatus(commentID uint, req *api.CommentUpdateRequest) (*model.Comment, error) {
	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return nil, errors.New("评论不存在")
	}

	comment.Status = req.Status

	if err := s.commentRepo.UpdateComment(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

// DeleteComment 删除评论
// @commentID 评论ID
// @return 错误
func (s *CommentService) DeleteComment(commentID uint) error {
	_, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}
	return s.commentRepo.DeleteComment(commentID)
}

package services

import (
	"errors"
	"klog-backend/internal/api"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"
	"klog-backend/internal/utils"
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

// toCommentResponse 将Comment模型转换为CommentResponse（带IP脱敏）
// @comment Comment模型
// @return CommentResponse
func toCommentResponse(comment *model.Comment) api.CommentResponse {
	response := api.CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Name:      comment.Name,
		Email:     comment.Email,
		Content:   comment.Content,
		IP:        utils.MaskIP(comment.IP), // IP脱敏
		Status:    comment.Status,
		ParentID:  comment.ParentID,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 递归处理回复列表
	if len(comment.Replies) > 0 {
		response.Replies = make([]api.CommentResponse, len(comment.Replies))
		for i, reply := range comment.Replies {
			response.Replies[i] = toCommentResponse(&reply)
		}
	}

	return response
}

// CreateComment 创建评论
// @postID 文章ID
// @req 创建评论请求
// @userID 用户ID（可为nil，表示游客）
// @ip IP地址
// @return 评论响应（IP已脱敏）, 错误
func (s *CommentService) CreateComment(postID uint, req *api.CommentCreateRequest, userID *uint, ip string) (*api.CommentResponse, error) {
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

	// 转换为响应格式（IP脱敏）
	response := toCommentResponse(comment)
	return &response, nil
}

// GetCommentsByPostID 根据文章ID获取评论列表
// @postID 文章ID
// @return 评论响应列表（IP已脱敏）, 错误
func (s *CommentService) GetCommentsByPostID(postID uint) ([]api.CommentResponse, error) {
	comments, err := s.commentRepo.GetCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式（IP脱敏）
	responses := make([]api.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = toCommentResponse(&comment)
	}

	return responses, nil
}

// UpdateCommentStatus 更新评论状态
// @commentID 评论ID
// @req 更新评论请求
// @return 评论响应（IP已脱敏）, 错误
func (s *CommentService) UpdateCommentStatus(commentID uint, req *api.CommentUpdateRequest) (*api.CommentResponse, error) {
	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return nil, errors.New("评论不存在")
	}

	comment.Status = req.Status

	if err := s.commentRepo.UpdateComment(comment); err != nil {
		return nil, err
	}

	// 转换为响应格式（IP脱敏）
	response := toCommentResponse(comment)
	return &response, nil
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

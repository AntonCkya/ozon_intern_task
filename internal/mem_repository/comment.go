package mem_repository

import (
	"context"
	"errors"
	"sync"

	"github.com/AntonCkya/ozon_habr/internal/repo_models"
)

type CommentRepository struct {
	mu       sync.RWMutex
	comments map[int]*repo_models.Comment
	nextID   int
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		comments: make(map[int]*repo_models.Comment),
		nextID:   1,
	}
}

func (r *CommentRepository) CreateComment(ctx context.Context, content string, userID, postID, parentID int) (*repo_models.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	comment := &repo_models.Comment{
		ID:       r.nextID,
		Content:  content,
		UserID:   userID,
		PostID:   postID,
		ParentID: nil,
	}

	if parentID != -1 {
		comment.ParentID = &parentID
	}

	r.comments[comment.ID] = comment
	r.nextID++

	return &repo_models.Comment{
		ID:       comment.ID,
		Content:  comment.Content,
		UserID:   comment.UserID,
		PostID:   comment.PostID,
		ParentID: comment.ParentID,
	}, nil
}

func (r *CommentRepository) GetCommentsByPostID(ctx context.Context, postID int, limit int, offset int) ([]*repo_models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var postComments []*repo_models.Comment
	for _, comment := range r.comments {
		if comment.PostID == postID {
			postComments = append(postComments, comment)
		}
	}
	start := offset
	if start > len(postComments) {
		start = len(postComments)
	}
	end := start + limit
	if end > len(postComments) {
		end = len(postComments)
	}
	comments := postComments[start:end]

	result := make([]*repo_models.Comment, 0, len(comments))
	for _, comment := range comments {
		result = append(result, &repo_models.Comment{
			ID:       comment.ID,
			Content:  comment.Content,
			UserID:   comment.UserID,
			PostID:   comment.PostID,
			ParentID: comment.ParentID,
		})
	}

	return result, nil
}

func (r *CommentRepository) GetCommentsByPostIDs(ctx context.Context, postIDs []int) ([]*repo_models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	postIDSet := make(map[int]struct{})
	for _, id := range postIDs {
		postIDSet[id] = struct{}{}
	}
	var result []*repo_models.Comment
	for _, comment := range r.comments {
		if _, exists := postIDSet[comment.PostID]; exists {
			result = append(result, &repo_models.Comment{
				ID:       comment.ID,
				Content:  comment.Content,
				UserID:   comment.UserID,
				PostID:   comment.PostID,
				ParentID: comment.ParentID,
			})
		}
	}

	return result, nil
}

func (r *CommentRepository) GetReplies(ctx context.Context, parentID int) ([]*repo_models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var replies []*repo_models.Comment
	for _, comment := range r.comments {
		if comment.ParentID != nil && *comment.ParentID == parentID {
			replies = append(replies, &repo_models.Comment{
				ID:       comment.ID,
				Content:  comment.Content,
				UserID:   comment.UserID,
				PostID:   comment.PostID,
				ParentID: comment.ParentID,
			})
		}
	}

	return replies, nil
}

func (r *CommentRepository) GetCommentByID(ctx context.Context, id int) (*repo_models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[id]
	if !exists {
		return nil, errors.New("comment not found")
	}

	return &repo_models.Comment{
		ID:       comment.ID,
		Content:  comment.Content,
		UserID:   comment.UserID,
		PostID:   comment.PostID,
		ParentID: comment.ParentID,
	}, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, id int, content string) (*repo_models.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	comment, exists := r.comments[id]
	if !exists {
		return nil, errors.New("comment not found")
	}

	comment.Content = content

	return &repo_models.Comment{
		ID:       comment.ID,
		Content:  comment.Content,
		UserID:   comment.UserID,
		PostID:   comment.PostID,
		ParentID: comment.ParentID,
	}, nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.comments[id]
	if !exists {
		return errors.New("comment not found")
	}

	delete(r.comments, id)
	return nil
}

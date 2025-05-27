package mem_repository

import (
	"context"
	"errors"
	"sync"

	"github.com/AntonCkya/ozon_habr/internal/repo_models"
)

type PostRepository struct {
	mu     sync.RWMutex
	posts  map[int]*repo_models.Post
	nextID int
}

func NewPostRepository() *PostRepository {
	return &PostRepository{
		posts:  make(map[int]*repo_models.Post),
		nextID: 1,
	}
}

func (r *PostRepository) CreatePost(ctx context.Context, title, content string, userID int, commentable bool) (*repo_models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	post := &repo_models.Post{
		ID:          r.nextID,
		Title:       title,
		Content:     content,
		UserID:      userID,
		Commentable: commentable,
	}

	r.posts[post.ID] = post
	r.nextID++

	return &repo_models.Post{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		UserID:      post.UserID,
		Commentable: post.Commentable,
	}, nil
}

func (r *PostRepository) GetPostByID(ctx context.Context, id int) (*repo_models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}

	return &repo_models.Post{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		UserID:      post.UserID,
		Commentable: post.Commentable,
	}, nil
}

func (r *PostRepository) GetPosts(ctx context.Context, limit int, offset int) ([]*repo_models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allPosts := make([]*repo_models.Post, 0, len(r.posts))
	for _, post := range r.posts {
		allPosts = append(allPosts, post)
	}
	start := offset
	if start > len(allPosts) {
		start = len(allPosts)
	}
	end := start + limit
	if end > len(allPosts) {
		end = len(allPosts)
	}
	posts := allPosts[start:end]

	result := make([]*repo_models.Post, 0, len(posts))
	for _, post := range posts {
		result = append(result, &repo_models.Post{
			ID:          post.ID,
			Title:       post.Title,
			Content:     post.Content,
			UserID:      post.UserID,
			Commentable: post.Commentable,
		})
	}

	return result, nil
}

func (r *PostRepository) GetPostsByUserId(ctx context.Context, limit int, offset int, userId int) ([]*repo_models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userPosts []*repo_models.Post
	for _, post := range r.posts {
		if post.UserID == userId {
			userPosts = append(userPosts, post)
		}
	}
	start := offset
	if start > len(userPosts) {
		start = len(userPosts)
	}
	end := start + limit
	if end > len(userPosts) {
		end = len(userPosts)
	}
	posts := userPosts[start:end]

	result := make([]*repo_models.Post, 0, len(posts))
	for _, post := range posts {
		result = append(result, &repo_models.Post{
			ID:          post.ID,
			Title:       post.Title,
			Content:     post.Content,
			UserID:      post.UserID,
			Commentable: post.Commentable,
		})
	}

	return result, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, id int, title, content string, userID int, commentable bool) (*repo_models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	post, exists := r.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}

	post.Title = title
	post.Content = content
	post.Commentable = commentable

	return &repo_models.Post{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		UserID:      post.UserID,
		Commentable: post.Commentable,
	}, nil
}

func (r *PostRepository) DeletePost(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.posts[id]
	if !exists {
		return errors.New("post not found")
	}

	delete(r.posts, id)
	return nil
}

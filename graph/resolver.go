package graph

import (
	"context"
	"database/sql"

	"github.com/AntonCkya/ozon_habr/internal/mem_repository"
	"github.com/AntonCkya/ozon_habr/internal/pg_repository"
	"github.com/AntonCkya/ozon_habr/internal/repo_models"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type UserRepoInterface interface {
	CreateUser(ctx context.Context, username string, password string) (*repo_models.User, error)
	GetUserByID(ctx context.Context, id int) (*repo_models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*repo_models.User, error)
	GetUsersByIDs(ctx context.Context, ids []int) ([]*repo_models.User, error)
}

type PostRepoInterface interface {
	CreatePost(ctx context.Context, title string, content string, userID int, commentable bool) (*repo_models.Post, error)
	DeletePost(ctx context.Context, id int) error
	GetPostByID(ctx context.Context, id int) (*repo_models.Post, error)
	GetPosts(ctx context.Context, limit int, offset int) ([]*repo_models.Post, error)
	GetPostsByUserId(ctx context.Context, limit int, offset int, userId int) ([]*repo_models.Post, error)
	UpdatePost(ctx context.Context, id int, title string, content string, userID int, commentable bool) (*repo_models.Post, error)
}

type CommentRepoInterface interface {
	CreateComment(ctx context.Context, content string, userID int, postID int, parentID int) (*repo_models.Comment, error)
	DeleteComment(ctx context.Context, id int) error
	GetCommentByID(ctx context.Context, id int) (*repo_models.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID int, limit int, offset int) ([]*repo_models.Comment, error)
	GetCommentsByPostIDs(ctx context.Context, postIDs []int) ([]*repo_models.Comment, error)
	GetReplies(ctx context.Context, parentID int) ([]*repo_models.Comment, error)
	UpdateComment(ctx context.Context, id int, content string) (*repo_models.Comment, error)
}

type Resolver struct {
	UserRepo    UserRepoInterface
	PostRepo    PostRepoInterface
	CommentRepo CommentRepoInterface
}

func NewPgResolver(db *sql.DB) *Resolver {
	return &Resolver{
		UserRepo:    pg_repository.NewUserRepository(db),
		PostRepo:    pg_repository.NewPostRepository(db),
		CommentRepo: pg_repository.NewCommentRepository(db),
	}
}

func NewMemResolver() *Resolver {
	return &Resolver{
		UserRepo:    mem_repository.NewUserRepository(),
		PostRepo:    mem_repository.NewPostRepository(),
		CommentRepo: mem_repository.NewCommentRepository(),
	}
}

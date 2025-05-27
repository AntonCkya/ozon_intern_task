package pg_repository

import (
	"context"
	"database/sql"

	"github.com/AntonCkya/ozon_habr/internal/repo_models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

const (
	CreatePostQuery = `
		INSERT INTO posts (title, content, user_id, commentable)
		VALUES ($1, $2, $3, $4)
	    RETURNING id, title, content, user_id, commentable;
	`
	GetPostByIdQuery = `
		SELECT id, title, content, user_id, commentable
		FROM posts
		WHERE id = $1;
	`
	GetPostsByUserIdQuery = `
		SELECT id, title, content, user_id, commentable
		FROM posts
		WHERE user_id = $1
		LIMIT $2
		OFFSET $3;
	`
	GetPostsQuery = `
		SELECT id, title, content, user_id, commentable
		FROM posts
		LIMIT $1
		OFFSET $2;
	`
	UpdatePostQuery = `
		UPDATE posts
		SET
		title = $1,
		content = $2,
		commentable = $5
	    WHERE id = $3 AND user_id = $4
	    RETURNING id, title, content, user_id, commentable;
	`
	DeletePostQuery = `
		DELETE FROM posts
		WHERE id = $1;
	`
)

func (r *PostRepository) CreatePost(ctx context.Context, title, content string, userID int, commentable bool) (*repo_models.Post, error) {
	var post repo_models.Post

	row := r.db.QueryRowContext(ctx, CreatePostQuery, title, content, userID, commentable)
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.Commentable,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetPostByID(ctx context.Context, id int) (*repo_models.Post, error) {
	var post repo_models.Post
	row := r.db.QueryRowContext(ctx, GetPostByIdQuery, id)
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.Commentable,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetPosts(ctx context.Context, limit int, offset int) ([]*repo_models.Post, error) {
	rows, err := r.db.QueryContext(ctx, GetPostsQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*repo_models.Post
	for rows.Next() {
		var post repo_models.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.UserID,
			&post.Commentable,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) GetPostsByUserId(ctx context.Context, limit int, offset int, userId int) ([]*repo_models.Post, error) {
	rows, err := r.db.QueryContext(ctx, GetPostsByUserIdQuery, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*repo_models.Post
	for rows.Next() {
		var post repo_models.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.UserID,
			&post.Commentable,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, id int, title, content string, userID int, commentable bool) (*repo_models.Post, error) {
	var post repo_models.Post
	row := r.db.QueryRowContext(ctx, UpdatePostQuery, title, content, id, userID, commentable)
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.Commentable,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) DeletePost(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, DeletePostQuery, id)
	if err != nil {
		return err
	}

	return nil
}

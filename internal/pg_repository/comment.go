package pg_repository

import (
	"context"
	"database/sql"

	"github.com/AntonCkya/ozon_habr/internal/repo_models"
	"github.com/lib/pq"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

const (
	CreateCommentQuery = `
		INSERT INTO comments (content, user_id, post_id, parent_id) 
		VALUES ($1, $2, $3, $4) 
	    RETURNING id, content, user_id, post_id, parent_id;
	`
	GetCommentsByPostIdQuery = `
		SELECT id, content, user_id, post_id, parent_id
	    FROM comments
		WHERE post_id = $1
		LIMIT $2
		OFFSET $3;
	`
	GetCommentsByPostIdBulkQuery = `
		SELECT id, content, user_id, post_id, parent_id
	    FROM comments
		WHERE post_id = ANY($1);
	`
	GetRepliesQuery = `
		SELECT id, content, user_id, post_id, parent_id
	    FROM comments
		WHERE parent_id = $1;
	`
	UpdateCommentQuery = `
		UPDATE comments
		SET
		content = $2
	    WHERE id = $1
	    RETURNING id, content, user_id, post_id, parent_id;
	`
	DeleteCommentQuery = `
		DELETE FROM comments
		WHERE id = $1;
	`
	GetCommentQuery = `
		SELECT id, content, user_id, post_id, parent_id
	    FROM comments
		WHERE id = $1;
	`
)

func (r *CommentRepository) CreateComment(ctx context.Context, content string, userID, postID, parentID int) (*repo_models.Comment, error) {
	var comment repo_models.Comment
	var row *sql.Row
	if parentID == -1 {
		row = r.db.QueryRowContext(ctx, CreateCommentQuery, content, userID, postID, nil)
	} else {
		row = r.db.QueryRowContext(ctx, CreateCommentQuery, content, userID, postID, parentID)
	}
	err := row.Scan(
		&comment.ID,
		&comment.Content,
		&comment.UserID,
		&comment.PostID,
		&comment.ParentID,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) GetCommentsByPostID(ctx context.Context, postID int, limit int, offset int) ([]*repo_models.Comment, error) {
	rows, err := r.db.QueryContext(ctx, GetCommentsByPostIdQuery, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*repo_models.Comment
	for rows.Next() {
		var comment repo_models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.UserID,
			&comment.PostID,
			&comment.ParentID,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (r *CommentRepository) GetCommentsByPostIDs(ctx context.Context, postIDs []int) ([]*repo_models.Comment, error) {
	rows, err := r.db.QueryContext(ctx, GetCommentsByPostIdBulkQuery, pq.Array(postIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*repo_models.Comment
	for rows.Next() {
		var comment repo_models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.UserID,
			&comment.PostID,
			&comment.ParentID,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (r *CommentRepository) GetReplies(ctx context.Context, parentID int) ([]*repo_models.Comment, error) {
	rows, err := r.db.QueryContext(ctx, GetRepliesQuery, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*repo_models.Comment
	for rows.Next() {
		var comment repo_models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.UserID,
			&comment.PostID,
			&comment.ParentID,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (r *CommentRepository) GetCommentByID(ctx context.Context, id int) (*repo_models.Comment, error) {
	var comment repo_models.Comment
	row := r.db.QueryRowContext(ctx, GetCommentQuery, id)
	err := row.Scan(
		&comment.ID,
		&comment.Content,
		&comment.UserID,
		&comment.PostID,
		&comment.ParentID,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, id int, content string) (*repo_models.Comment, error) {
	var comment repo_models.Comment
	row := r.db.QueryRowContext(ctx, UpdateCommentQuery, id, content)
	err := row.Scan(
		&comment.ID,
		&comment.Content,
		&comment.UserID,
		&comment.PostID,
		&comment.ParentID,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, DeleteCommentQuery, id)
	if err != nil {
		return err
	}

	return nil
}

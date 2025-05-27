package pg_repository

import (
	"context"
	"database/sql"

	"github.com/AntonCkya/ozon_habr/internal/repo_models"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

const (
	CreateUserQuery = `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username;
	`
	GetUserByIdQuery = `
		SELECT id, username, password_hash
		FROM users 
		WHERE id = $1;
	`
	GetUserByNameQuery = `
		SELECT id, username, password_hash
		FROM users 
		WHERE username = $1;
	`
	// для решения N+1
	GetUsersByIdBulkQuery = `
        SELECT id, username, password_hash
        FROM users 
        WHERE id = ANY($1);
    `
)

func (r *UserRepository) CreateUser(ctx context.Context, username, password string) (*repo_models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user repo_models.User
	row := r.db.QueryRowContext(ctx, CreateUserQuery, username, string(hashedPassword))
	err = row.Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*repo_models.User, error) {
	user, err := r.GetUser(ctx, true, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*repo_models.User, error) {
	user, err := r.GetUser(ctx, false, username)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUser(ctx context.Context, by_id bool, payload any) (*repo_models.User, error) {
	var user repo_models.User

	var query string
	if by_id {
		query = GetUserByIdQuery
	} else {
		query = GetUserByNameQuery
	}

	row := r.db.QueryRowContext(ctx, query, payload)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUsersByIDs(ctx context.Context, ids []int) ([]*repo_models.User, error) {
	var users []*repo_models.User

	rows, err := r.db.QueryContext(ctx, GetUsersByIdBulkQuery, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user repo_models.User
		err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

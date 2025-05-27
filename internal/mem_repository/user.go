package mem_repository

import (
	"context"
	"errors"
	"sync"

	"github.com/AntonCkya/ozon_habr/internal/repo_models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	mu     sync.RWMutex
	users  map[int]*repo_models.User
	nextID int
}

// чтобы REST и GraphQL обработчики работали вместе при in memory режиме

var instance *UserRepository
var once sync.Once

func NewUserRepository() *UserRepository {
	once.Do(func() {
		instance = &UserRepository{
			users:  make(map[int]*repo_models.User),
			nextID: 1,
		}
	})
	return instance
}

func (r *UserRepository) CreateUser(ctx context.Context, username, password string) (*repo_models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range r.users {
		if user.Username == username {
			return nil, errors.New("user already exists")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &repo_models.User{
		ID:           r.nextID,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	r.users[user.ID] = user
	r.nextID++

	return &repo_models.User{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*repo_models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return &repo_models.User{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*repo_models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return &repo_models.User{
				ID:           user.ID,
				Username:     user.Username,
				PasswordHash: user.PasswordHash,
			}, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *UserRepository) GetUsersByIDs(ctx context.Context, ids []int) ([]*repo_models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*repo_models.User, 0, len(ids))
	for _, id := range ids {
		if user, exists := r.users[id]; exists {
			users = append(users, &repo_models.User{
				ID:           user.ID,
				Username:     user.Username,
				PasswordHash: user.PasswordHash,
			})
		}
	}

	return users, nil
}

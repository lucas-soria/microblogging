package users

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=users

var (
	ErrUserNotFound  = NewRepositoryError("user not found")
	ErrHandlerExists = NewRepositoryError("handler already exists")
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, handler string) (*User, error)
	DeleteUser(ctx context.Context, handler string) error
	FollowUser(ctx context.Context, followerHandler string, followeeHandler string) error
	UnfollowUser(ctx context.Context, followerHandler string, followeeHandler string) error
	GetUserFollowers(ctx context.Context, followeeHandler string) ([]User, error)
	GetUserFollowees(ctx context.Context, followerHandler string) ([]User, error)
}

type InMemoryUserRepository struct {
	mu     sync.RWMutex
	users  map[string]*User
	follow map[string]map[string]bool // followerHandler -> followeeHandler -> bool
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[string]*User),
		follow: make(map[string]map[string]bool),
	}
}

func (repository *InMemoryUserRepository) CreateUser(ctx context.Context, user *User) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	// Check if handler already exists
	for _, u := range repository.users {
		if u.Handler == user.Handler {
			return ErrHandlerExists
		}
	}

	// Generate new ID if not provided
	if user.Handler == "" {
		user.Handler = uuid.New().String()
	}

	repository.users[user.Handler] = user
	return nil
}

func (repository *InMemoryUserRepository) GetUser(ctx context.Context, handler string) (*User, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	user, exists := repository.users[handler]
	if !exists {
		return nil, ErrUserNotFound
	}

	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy, nil
}

func (repository *InMemoryUserRepository) DeleteUser(ctx context.Context, handler string) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	if _, exists := repository.users[handler]; !exists {
		return ErrUserNotFound
	}

	// Remove user from users map
	delete(repository.users, handler)

	// Remove user from follow relationships
	delete(repository.follow, handler) // Remove user's following relationships
	for followerID := range repository.follow {
		delete(repository.follow[followerID], handler) // Remove user from others' followers
	}

	return nil
}

func (repository *InMemoryUserRepository) FollowUser(ctx context.Context, followerHandler string, followeeHandler string) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	// Check if both users exist
	if _, exists := repository.users[followerHandler]; !exists {
		return ErrUserNotFound
	}
	if _, exists := repository.users[followeeHandler]; !exists {
		return ErrUserNotFound
	}

	// Initialize follower's follow map if it doesn't exist
	if repository.follow[followerHandler] == nil {
		repository.follow[followerHandler] = make(map[string]bool)
	}

	repository.follow[followerHandler][followeeHandler] = true
	return nil
}

func (repository *InMemoryUserRepository) UnfollowUser(ctx context.Context, followerHandler string, followeeHandler string) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	if repository.follow[followerHandler] != nil {
		delete(repository.follow[followerHandler], followeeHandler)
	}
	return nil
}

func (repository *InMemoryUserRepository) GetUserFollowers(ctx context.Context, handler string) ([]User, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	// Check if user exists
	if _, exists := repository.users[handler]; !exists {
		return nil, ErrUserNotFound
	}

	var followers []User
	for followerID, followees := range repository.follow {
		if followees[handler] {
			if user, exists := repository.users[followerID]; exists {
				followers = append(followers, *user)
			}
		}
	}

	return followers, nil
}

func (repository *InMemoryUserRepository) GetUserFollowees(ctx context.Context, handler string) ([]User, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	// Check if user exists
	if _, exists := repository.users[handler]; !exists {
		return nil, ErrUserNotFound
	}

	var following []User
	for followeeID := range repository.follow[handler] {
		if user, exists := repository.users[followeeID]; exists {
			following = append(following, *user)
		}
	}
	return following, nil
}

// Errors
type RepositoryError struct {
	message string
}

func (e *RepositoryError) Error() string {
	return e.message
}

func NewRepositoryError(message string) error {
	return &RepositoryError{message: message}
}

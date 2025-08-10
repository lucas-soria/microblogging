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

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, handler string) (*User, error)
	DeleteUser(ctx context.Context, handler string) error
	FollowUser(ctx context.Context, followRequest FollowRequest) error
	UnfollowUser(ctx context.Context, followRequest FollowRequest) error
	GetUserFollowers(ctx context.Context, handler string) ([]User, error)
	GetUserFollowing(ctx context.Context, handler string) ([]User, error)
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

func (r *InMemoryUserRepository) CreateUser(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if handler already exists
	for _, u := range r.users {
		if u.Handler == user.Handler {
			return ErrHandlerExists
		}
	}

	// Generate new ID if not provided
	if user.Handler == "" {
		user.Handler = uuid.New().String()
	}

	r.users[user.Handler] = user
	return nil
}

func (r *InMemoryUserRepository) GetUser(ctx context.Context, handler string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[handler]
	if !exists {
		return nil, ErrUserNotFound
	}

	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy, nil
}

func (r *InMemoryUserRepository) DeleteUser(ctx context.Context, handler string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[handler]; !exists {
		return ErrUserNotFound
	}

	// Remove user from users map
	delete(r.users, handler)

	// Remove user from follow relationships
	delete(r.follow, handler) // Remove user's following relationships
	for followerID := range r.follow {
		delete(r.follow[followerID], handler) // Remove user from others' followers
	}

	return nil
}

func (r *InMemoryUserRepository) FollowUser(ctx context.Context, followRequest FollowRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if both users exist
	if _, exists := r.users[followRequest.FollowerHandler]; !exists {
		return ErrUserNotFound
	}
	if _, exists := r.users[followRequest.FolloweeHandler]; !exists {
		return ErrUserNotFound
	}

	// Initialize follower's follow map if it doesn't exist
	if r.follow[followRequest.FollowerHandler] == nil {
		r.follow[followRequest.FollowerHandler] = make(map[string]bool)
	}

	r.follow[followRequest.FollowerHandler][followRequest.FolloweeHandler] = true
	return nil
}

func (r *InMemoryUserRepository) UnfollowUser(ctx context.Context, followRequest FollowRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.follow[followRequest.FollowerHandler] != nil {
		delete(r.follow[followRequest.FollowerHandler], followRequest.FolloweeHandler)
	}
	return nil
}

func (r *InMemoryUserRepository) GetUserFollowers(ctx context.Context, handler string) ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if user exists
	if _, exists := r.users[handler]; !exists {
		return nil, ErrUserNotFound
	}

	var followers []User
	for followerID, followees := range r.follow {
		if followees[handler] {
			if user, exists := r.users[followerID]; exists {
				followers = append(followers, *user)
			}
		}
	}

	return followers, nil
}

func (r *InMemoryUserRepository) GetUserFollowing(ctx context.Context, handler string) ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if user exists
	if _, exists := r.users[handler]; !exists {
		return nil, ErrUserNotFound
	}

	var following []User
	for followeeID := range r.follow[handler] {
		if user, exists := r.users[followeeID]; exists {
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

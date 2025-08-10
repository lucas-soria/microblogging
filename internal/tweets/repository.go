package tweets

import (
	"context"
	"sync"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=tweets

// Repository defines the interface for tweet data operations
type Repository interface {
	Create(ctx context.Context, tweet *Tweet) (*Tweet, error)
	GetByID(ctx context.Context, id string) (*Tweet, error)
	GetByUserID(ctx context.Context, userID string) ([]*Tweet, error)
	Delete(ctx context.Context, id string) error
}

// InMemoryTweetRepository is an in-memory implementation of the Repository interface
type InMemoryTweetRepository struct {
	tweets map[string]*Tweet
	mu     sync.RWMutex
}

// NewInMemoryTweetRepository creates a new in-memory tweet repository
func NewInMemoryTweetRepository() *InMemoryTweetRepository {
	return &InMemoryTweetRepository{
		tweets: make(map[string]*Tweet),
	}
}

func (r *InMemoryTweetRepository) Create(ctx context.Context, tweet *Tweet) (*Tweet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tweets[tweet.ID] = tweet
	return tweet, nil
}

func (r *InMemoryTweetRepository) GetByID(ctx context.Context, id string) (*Tweet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tweet, exists := r.tweets[id]
	if !exists {
		return nil, nil
	}
	return tweet, nil
}

func (r *InMemoryTweetRepository) GetByUserID(ctx context.Context, userID string) ([]*Tweet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userTweets []*Tweet
	for _, tweet := range r.tweets {
		if tweet.Handler == userID {
			userTweets = append(userTweets, tweet)
		}
	}
	return userTweets, nil
}

func (r *InMemoryTweetRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tweets, id)
	return nil
}

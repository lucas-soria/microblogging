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

func (repository *InMemoryTweetRepository) Create(ctx context.Context, tweet *Tweet) (*Tweet, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	repository.tweets[tweet.ID] = tweet
	return tweet, nil
}

func (repository *InMemoryTweetRepository) GetByID(ctx context.Context, id string) (*Tweet, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	tweet, exists := repository.tweets[id]
	if !exists {
		return nil, nil
	}
	return tweet, nil
}

func (repository *InMemoryTweetRepository) GetByUserID(ctx context.Context, userID string) ([]*Tweet, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	var userTweets []*Tweet
	for _, tweet := range repository.tweets {
		if tweet.Handler == userID {
			userTweets = append(userTweets, tweet)
		}
	}
	return userTweets, nil
}

func (repository *InMemoryTweetRepository) Delete(ctx context.Context, id string) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	delete(repository.tweets, id)
	return nil
}

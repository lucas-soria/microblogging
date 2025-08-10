package feed

import (
	"context"
	"sort"
	"sync"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=feed

// Repository defines the interface for feed data operations
type Repository interface {
	GetUserTimeline(ctx context.Context, userID string, limit, offset int) ([]*Tweet, error)
}

// InMemoryFeedRepository is an in-memory implementation of the Repository interface
type InMemoryFeedRepository struct {
	mu     sync.RWMutex
	tweets map[string][]*Tweet // userID -> []*Tweet
}

// NewInMemoryFeedRepository creates a new in-memory feed repository
func NewInMemoryFeedRepository() *InMemoryFeedRepository {
	return &InMemoryFeedRepository{
		tweets: make(map[string][]*Tweet),
	}
}

// GetUserTimeline retrieves the timeline for a user with pagination
func (repository *InMemoryFeedRepository) GetUserTimeline(ctx context.Context, userID string, limit, offset int) ([]*Tweet, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	timeline, exists := repository.tweets[userID]
	if !exists {
		return []*Tweet{}, nil
	}

	// Sort tweets by creation time (newest first)
	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].CreatedAt.After(timeline[j].CreatedAt)
	})

	// Apply pagination
	start := offset
	if start > len(timeline) {
		return []*Tweet{}, nil
	}

	end := offset + limit
	if end > len(timeline) {
		end = len(timeline)
	}

	result := make([]*Tweet, end-start)
	copy(result, timeline[start:end])

	return result, nil
}

// AddTweet adds a tweet to the feed of followers (helper method for testing)
func (repository *InMemoryFeedRepository) AddTweet(userID string, tweet *Tweet) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	repository.tweets[userID] = append(repository.tweets[userID], tweet)
}

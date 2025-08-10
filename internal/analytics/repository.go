package analytics

import (
	"context"
	"errors"
	"sync"
	"time"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=analytics

// Repository defines the interface for analytics data operations
type Repository interface {
	// User Analytics
	GetUserAnalytics(ctx context.Context, userID string) (*UserAnalytics, error)
	GetAllUserAnalytics(ctx context.Context) ([]*UserAnalytics, error)
	DeleteUserAnalytics(ctx context.Context, userID string) error

	// Event Processing
	ProcessEvent(ctx context.Context, event *Event) error
}

// InMemoryRepository is an in-memory implementation of the Repository interface
type InMemoryRepository struct {
	mu        sync.RWMutex
	analytics map[string]*UserAnalytics
	eventsMu  sync.RWMutex
	events    []*Event
}

// NewInMemoryRepository creates a new in-memory analytics repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		mu:        sync.RWMutex{},
		eventsMu:  sync.RWMutex{},
		analytics: map[string]*UserAnalytics{},
		events:    []*Event{},
	}
}

// GetUserAnalytics retrieves analytics for a specific user
func (repository *InMemoryRepository) GetUserAnalytics(ctx context.Context, userID string) (*UserAnalytics, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	analytics, exists := repository.analytics[userID]
	if !exists {
		return nil, errors.New("user analytics not found")
	}

	// Return a copy to prevent external modifications
	result := *analytics
	return &result, nil
}

// GetAllUserAnalytics retrieves analytics for all users
func (repository *InMemoryRepository) GetAllUserAnalytics(ctx context.Context) ([]*UserAnalytics, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	result := make([]*UserAnalytics, 0, len(repository.analytics))
	for _, analytics := range repository.analytics {
		// Create a copy to prevent external modifications
		analyticsCopy := *analytics
		result = append(result, &analyticsCopy)
	}

	return result, nil
}

// DeleteUserAnalytics deletes analytics data for a specific user
func (repository *InMemoryRepository) DeleteUserAnalytics(ctx context.Context, userID string) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	if _, exists := repository.analytics[userID]; !exists {
		return errors.New("user analytics not found")
	}

	delete(repository.analytics, userID)
	return nil
}

// ProcessEvent processes an analytics event
func (repository *InMemoryRepository) ProcessEvent(ctx context.Context, event *Event) error {
	repository.eventsMu.Lock()
	repository.events = append(repository.events, event)
	repository.eventsMu.Unlock()

	repository.mu.Lock()
	defer repository.mu.Unlock()

	now := time.Now()

	// Get or create user analytics
	analytics, exists := repository.analytics[event.Handler]
	if !exists {
		analytics = &UserAnalytics{
			Handler:   event.Handler,
			CreatedAt: now,
		}
	}

	// Update analytics based on event type
	switch event.EventType {
	case "tweet_created":
		// Mark user as active
		analytics.IsActive = true
		// If user has created many tweets, they might be an influencer
		// This is a simple heuristic - in a real app, we'd have more sophisticated logic
		tweetCount := 0
		repository.eventsMu.RLock()
		for _, e := range repository.events {
			if e.Handler == event.Handler && e.EventType == "tweet_created" {
				tweetCount++
			}
		}
		repository.eventsMu.RUnlock()

		if tweetCount > 100 { // Arbitrary threshold for demo
			analytics.IsInfluencer = true
		}

	case "timeline_viewed":
		// Mark user as active
		analytics.IsActive = true
	}

	analytics.UpdatedAt = now
	repository.analytics[event.Handler] = analytics

	return nil
}

package feed

import (
	"context"
	"errors"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=feed

// Service defines the business logic for feed operations
type Service interface {
	GetUserTimeline(ctx context.Context, userID string, limit, offset int) (*TimelineResponse, error)
}

type feedService struct {
	repo Repository
}

// NewFeedService creates a new feed service
func NewFeedService(repo Repository) Service {
	return &feedService{
		repo: repo,
	}
}

// GetUserTimeline retrieves the timeline for a user
func (s *feedService) GetUserTimeline(ctx context.Context, userID string, limit, offset int) (*TimelineResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Set default values if not provided
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	tweets, err := s.repo.GetUserTimeline(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Calculate next offset
	var nextOffset int
	if len(tweets) < limit {
		nextOffset = len(tweets)
	} else {
		nextOffset = limit
	}

	return &TimelineResponse{
		Tweets:     tweets,
		NextOffset: nextOffset,
	}, nil
}

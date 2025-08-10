package analytics

import (
	"context"
	"errors"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=analytics

// Service defines the business logic for analytics operations
type Service interface {
	// User Analytics
	GetUserAnalytics(ctx context.Context, userID string) (*UserAnalytics, error)
	GetAllUserAnalytics(ctx context.Context) ([]*UserAnalytics, error)
	DeleteUserAnalytics(ctx context.Context, userID string) error

	// Event Processing
	ProcessEvent(ctx context.Context, event *Event) error
}

type analyticsService struct {
	repo Repository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(repo Repository) Service {
	return &analyticsService{
		repo: repo,
	}
}

// GetUserAnalytics retrieves analytics for a specific user
func (s *analyticsService) GetUserAnalytics(ctx context.Context, userID string) (*UserAnalytics, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.repo.GetUserAnalytics(ctx, userID)
}

// GetAllUserAnalytics retrieves analytics for all users
func (s *analyticsService) GetAllUserAnalytics(ctx context.Context) ([]*UserAnalytics, error) {
	return s.repo.GetAllUserAnalytics(ctx)
}

// DeleteUserAnalytics deletes analytics data for a specific user
func (s *analyticsService) DeleteUserAnalytics(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	return s.repo.DeleteUserAnalytics(ctx, userID)
}

// ProcessEvent processes an analytics event
func (s *analyticsService) ProcessEvent(ctx context.Context, event *Event) error {
	if event == nil {
		return errors.New("event cannot be nil")
	}

	if event.Handler == "" {
		return errors.New("user ID is required in event")
	}

	if event.EventType == "" {
		return errors.New("event type is required")
	}

	return s.repo.ProcessEvent(ctx, event)
}

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

type service struct {
	repository Repository
}

// NewService creates a new analytics service
func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

// GetUserAnalytics retrieves analytics for a specific user
func (service *service) GetUserAnalytics(ctx context.Context, userID string) (*UserAnalytics, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return service.repository.GetUserAnalytics(ctx, userID)
}

// GetAllUserAnalytics retrieves analytics for all users
func (service *service) GetAllUserAnalytics(ctx context.Context) ([]*UserAnalytics, error) {
	return service.repository.GetAllUserAnalytics(ctx)
}

// DeleteUserAnalytics deletes analytics data for a specific user
func (service *service) DeleteUserAnalytics(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	return service.repository.DeleteUserAnalytics(ctx, userID)
}

// ProcessEvent processes an analytics event
func (service *service) ProcessEvent(ctx context.Context, event *Event) error {
	if event == nil {
		return errors.New("event cannot be nil")
	}

	if event.Handler == "" {
		return errors.New("user ID is required in event")
	}

	if event.EventType == "" {
		return errors.New("event type is required")
	}

	return service.repository.ProcessEvent(ctx, event)
}

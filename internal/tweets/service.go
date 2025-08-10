package tweets

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=tweets

// Service defines the business logic for tweet operations
type Service interface {
	CreateTweet(ctx context.Context, tweetToCreate *Tweet) (*Tweet, error)
	GetTweet(ctx context.Context, id string) (*Tweet, error)
	GetUserTweets(ctx context.Context, userID string) ([]*Tweet, error)
	DeleteTweet(ctx context.Context, id string) error
}

type service struct {
	repository Repository
}

// NewService creates a new tweet service
func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (service *service) CreateTweet(ctx context.Context, tweetToCreate *Tweet) (*Tweet, error) {
	if tweetToCreate.Content.Text == "" {
		return nil, errors.New("tweet content cannot be empty")
	}

	tweetToCreate.ID = uuid.New().String()
	tweetToCreate.CreatedAt = time.Now().UTC()

	return service.repository.Create(ctx, tweetToCreate)
}

func (service *service) GetTweet(ctx context.Context, id string) (*Tweet, error) {
	if id == "" {
		return nil, errors.New("tweet ID cannot be empty")
	}

	return service.repository.GetByID(ctx, id)
}

func (service *service) GetUserTweets(ctx context.Context, userID string) ([]*Tweet, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	return service.repository.GetByUserID(ctx, userID)
}

func (service *service) DeleteTweet(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("tweet ID cannot be empty")
	}

	// Check if tweet exists
	foundTweet, err := service.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if foundTweet == nil {
		return errors.New("tweet not found")
	}

	return service.repository.Delete(ctx, id)
}

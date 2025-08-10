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
	CreateTweet(ctx context.Context, req *CreateTweetRequest) (*Tweet, error)
	GetTweet(ctx context.Context, id string) (*Tweet, error)
	GetUserTweets(ctx context.Context, userID string) ([]*Tweet, error)
	DeleteTweet(ctx context.Context, id string) error
}

type tweetService struct {
	repo Repository
}

// NewTweetService creates a new tweet service
func NewTweetService(repo Repository) Service {
	return &tweetService{
		repo: repo,
	}
}

func (s *tweetService) CreateTweet(ctx context.Context, req *CreateTweetRequest) (*Tweet, error) {
	if req.Content.Text == "" {
		return nil, errors.New("tweet content cannot be empty")
	}

	tweet := &Tweet{
		ID:        uuid.New().String(),
		Handler:   req.Handler,
		Content:   req.Content,
		CreatedAt: time.Now().UTC(),
	}

	return s.repo.Create(ctx, tweet)
}

func (s *tweetService) GetTweet(ctx context.Context, id string) (*Tweet, error) {
	if id == "" {
		return nil, errors.New("tweet ID cannot be empty")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *tweetService) GetUserTweets(ctx context.Context, userID string) ([]*Tweet, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	return s.repo.GetByUserID(ctx, userID)
}

func (s *tweetService) DeleteTweet(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("tweet ID cannot be empty")
	}

	// Check if tweet exists
	foundTweet, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if foundTweet == nil {
		return errors.New("tweet not found")
	}

	return s.repo.Delete(ctx, id)
}

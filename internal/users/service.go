package users

import (
	"context"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=users

type UserService interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	FollowUser(ctx context.Context, req FollowRequest) error
	UnfollowUser(ctx context.Context, req FollowRequest) error
	GetUserFollowers(ctx context.Context, userID string) ([]User, error)
	GetUserFollowees(ctx context.Context, userID string) ([]User, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *User) (*User, error) {
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *userService) FollowUser(ctx context.Context, req FollowRequest) error {
	return s.repo.FollowUser(ctx, req)
}

func (s *userService) UnfollowUser(ctx context.Context, req FollowRequest) error {
	return s.repo.UnfollowUser(ctx, req)
}

func (s *userService) GetUserFollowers(ctx context.Context, userID string) ([]User, error) {
	return s.repo.GetUserFollowers(ctx, userID)
}

func (s *userService) GetUserFollowees(ctx context.Context, userID string) ([]User, error) {
	return s.repo.GetUserFollowing(ctx, userID)
}

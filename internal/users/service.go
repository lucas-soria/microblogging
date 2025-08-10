package users

import (
	"context"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=users

type Service interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	FollowUser(ctx context.Context, followerHandler string, followeeHandler string) error
	UnfollowUser(ctx context.Context, followerHandler string, followeeHandler string) error
	GetUserFollowers(ctx context.Context, followeeHandler string) ([]User, error)
	GetUserFollowees(ctx context.Context, followerHandler string) ([]User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (service *service) CreateUser(ctx context.Context, user *User) (*User, error) {
	if err := service.repository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (service *service) GetUser(ctx context.Context, id string) (*User, error) {
	return service.repository.GetUser(ctx, id)
}

func (service *service) DeleteUser(ctx context.Context, id string) error {
	return service.repository.DeleteUser(ctx, id)
}

func (service *service) FollowUser(ctx context.Context, followerHandler string, followeeHandler string) error {
	return service.repository.FollowUser(ctx, followerHandler, followeeHandler)
}

func (service *service) UnfollowUser(ctx context.Context, followerHandler string, followeeHandler string) error {
	return service.repository.UnfollowUser(ctx, followerHandler, followeeHandler)
}

func (service *service) GetUserFollowers(ctx context.Context, followeeHandler string) ([]User, error) {
	return service.repository.GetUserFollowers(ctx, followeeHandler)
}

func (service *service) GetUserFollowees(ctx context.Context, followerHandler string) ([]User, error) {
	return service.repository.GetUserFollowees(ctx, followerHandler)
}

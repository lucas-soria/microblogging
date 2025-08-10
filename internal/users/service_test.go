package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := &User{
		Handler:   "testuser",
		FirstName: "test",
		LastName:  "user",
	}

	type want struct {
		user *User
		err  error
	}

	tt := []struct {
		name         string
		expectations func(user *User)
		want         want
	}{
		{
			name: "successful user creation",
			expectations: func(user *User) {
				mockRepo.EXPECT().CreateUser(ctx, user).
					Return(nil).
					Times(1)
			},
			want: want{
				user: user,
				err:  nil,
			},
		},
		{
			name: "failed user creation",
			expectations: func(user *User) {
				mockRepo.EXPECT().CreateUser(ctx, user).
					Return(ErrHandlerExists).
					Times(1)
			},
			want: want{
				user: nil,
				err:  ErrHandlerExists,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(user)

			user, err := service.CreateUser(ctx, user)
			assert.Equal(t, tc.want.err, err)
			assert.Equal(t, tc.want.user, user)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := &User{
		Handler:   "testuser",
		FirstName: "test",
		LastName:  "user",
	}

	type want struct {
		user *User
		err  error
	}

	tt := []struct {
		name         string
		expectations func(user *User)
		want         want
	}{
		{
			name: "successful user retrieval",
			expectations: func(user *User) {
				mockRepo.EXPECT().GetUser(ctx, user.Handler).
					Return(user, nil).
					Times(1)
			},
			want: want{
				user: user,
				err:  nil,
			},
		},
		{
			name: "failed user retrieval",
			expectations: func(user *User) {
				mockRepo.EXPECT().GetUser(ctx, user.Handler).
					Return(nil, ErrUserNotFound).
					Times(1)
			},
			want: want{
				user: nil,
				err:  ErrUserNotFound,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(user)

			user, err := service.GetUser(ctx, user.Handler)
			assert.Equal(t, tc.want.err, err)
			assert.Equal(t, tc.want.user, user)
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	userID := "testid"

	type want struct {
		err error
	}

	tt := []struct {
		name         string
		expectations func()
		want         want
	}{
		{
			name: "successful user deletion",
			expectations: func() {
				mockRepo.EXPECT().DeleteUser(ctx, userID).
					Return(nil).
					Times(1)
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "failed user deletion - not found",
			expectations: func() {
				mockRepo.EXPECT().DeleteUser(ctx, userID).
					Return(ErrUserNotFound).
					Times(1)
			},
			want: want{
				err: ErrUserNotFound,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			err := service.DeleteUser(ctx, userID)
			assert.Equal(t, tc.want.err, err)
		})
	}
}

func TestUserService_FollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	req := FollowRequest{
		FollowerHandler: "follower1",
		FolloweeHandler: "followee1",
	}

	type want struct {
		err error
	}

	tt := []struct {
		name         string
		expectations func()
		want         want
	}{
		{
			name: "successful follow",
			expectations: func() {
				mockRepo.EXPECT().FollowUser(ctx, req).
					Return(nil).
					Times(1)
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "failed follow - user not found",
			expectations: func() {
				mockRepo.EXPECT().FollowUser(ctx, req).
					Return(ErrUserNotFound).
					Times(1)
			},
			want: want{
				err: ErrUserNotFound,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			err := service.FollowUser(ctx, req)
			assert.Equal(t, tc.want.err, err)
		})
	}
}

func TestUserService_UnfollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	req := FollowRequest{
		FollowerHandler: "follower1",
		FolloweeHandler: "followee1",
	}

	type want struct {
		err error
	}

	tt := []struct {
		name         string
		expectations func()
		want         want
	}{
		{
			name: "successful unfollow",
			expectations: func() {
				mockRepo.EXPECT().UnfollowUser(ctx, req).
					Return(nil).
					Times(1)
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "failed unfollow - not following",
			expectations: func() {
				mockRepo.EXPECT().UnfollowUser(ctx, req).
					Return(ErrUserNotFound).
					Times(1)
			},
			want: want{
				err: ErrUserNotFound,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			err := service.UnfollowUser(ctx, req)
			assert.Equal(t, tc.want.err, err)
		})
	}
}

func TestUserService_GetUserFollowers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	userID := "testuser"
	followers := []User{
		{Handler: "follower1"},
		{Handler: "follower2"},
	}

	type want struct {
		followers []User
		err       error
	}

	tt := []struct {
		name         string
		expectations func()
		want         want
	}{
		{
			name: "successful get followers",
			expectations: func() {
				mockRepo.EXPECT().GetUserFollowers(ctx, userID).
					Return(followers, nil).
					Times(1)
			},
			want: want{
				followers: followers,
				err:       nil,
			},
		},
		{
			name: "no followers found",
			expectations: func() {
				mockRepo.EXPECT().GetUserFollowers(ctx, userID).
					Return(nil, nil).
					Times(1)
			},
			want: want{
				followers: nil,
				err:       nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			result, err := service.GetUserFollowers(ctx, userID)
			assert.Equal(t, tc.want.err, err)
			assert.Equal(t, tc.want.followers, result)
		})
	}
}

func TestUserService_GetUserFollowees(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	userID := "testuser"
	followees := []User{
		{Handler: "user1"},
		{Handler: "user2"},
	}

	type want struct {
		followees []User
		err       error
	}

	tt := []struct {
		name         string
		expectations func()
		want         want
	}{
		{
			name: "successful get followees",
			expectations: func() {
				mockRepo.EXPECT().GetUserFollowing(ctx, userID).
					Return(followees, nil).
					Times(1)
			},
			want: want{
				followees: followees,
				err:       nil,
			},
		},
		{
			name: "no followees found",
			expectations: func() {
				mockRepo.EXPECT().GetUserFollowing(ctx, userID).
					Return(nil, nil).
					Times(1)
			},
			want: want{
				followees: nil,
				err:       nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			result, err := service.GetUserFollowees(ctx, userID)
			assert.Equal(t, tc.want.err, err)
			assert.Equal(t, tc.want.followees, result)
		})
	}
}

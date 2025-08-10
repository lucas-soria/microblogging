package analytics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserAnalytics(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := NewMockRepository(ctrl)
	service := NewAnalyticsService(repoMock)

	type want struct {
		err       error
		analytics *UserAnalytics
	}

	tt := []struct {
		name         string
		expectations func(userID string)
		userID       string
		want         want
	}{
		{
			name: "success",
			expectations: func(userID string) {
				repoMock.EXPECT().
					GetUserAnalytics(gomock.Any(), userID).
					Return(&UserAnalytics{
						Handler:      userID,
						IsInfluencer: true,
						IsActive:     true,
					}, nil)
			},
			userID: "test-user-1",
			want: want{
				err: nil,
				analytics: &UserAnalytics{
					Handler:      "test-user-1",
					IsInfluencer: true,
					IsActive:     true,
				},
			},
		},
		{
			name:         "empty user ID",
			expectations: func(userID string) {},
			userID:       "",
			want: want{
				err: errors.New("user ID is required"),
			},
		},
		{
			name: "repository error",
			expectations: func(userID string) {
				repoMock.EXPECT().
					GetUserAnalytics(gomock.Any(), userID).
					Return(nil, errors.New("database error"))
			},
			userID: "test-user-1",
			want: want{
				err: errors.New("database error"),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.userID)

			result, err := service.GetUserAnalytics(ctx, tc.userID)

			if tc.want.err != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.want.err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want.analytics, result)
			}
		})
	}
}

func TestGetAllUserAnalytics(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := NewMockRepository(ctrl)
	service := NewAnalyticsService(repoMock)

	type want struct {
		analytics []*UserAnalytics
		err       error
	}

	tt := []struct {
		name         string
		expectations func()
		want         want
	}{
		{
			name: "success",
			expectations: func() {
				repoMock.EXPECT().
					GetAllUserAnalytics(gomock.Any()).
					Return([]*UserAnalytics{
						{Handler: "user1", IsActive: true},
						{Handler: "user2", IsInfluencer: true},
					}, nil)
			},
			want: want{
				analytics: []*UserAnalytics{
					{Handler: "user1", IsActive: true},
					{Handler: "user2", IsInfluencer: true},
				},
				err: nil,
			},
		},
		{
			name: "repository error",
			expectations: func() {
				repoMock.EXPECT().
					GetAllUserAnalytics(gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			want: want{
				err: errors.New("database error"),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			result, err := service.GetAllUserAnalytics(ctx)

			assert.Equal(t, tc.want.analytics, result)
			assert.Equal(t, tc.want.err, err)
		})
	}
}

func TestDeleteUserAnalytics(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := NewMockRepository(ctrl)
	service := NewAnalyticsService(repoMock)

	type want struct {
		err error
	}

	tt := []struct {
		name         string
		expectations func()
		userID       string
		want         want
	}{
		{
			name: "success",
			expectations: func() {
				repoMock.EXPECT().
					DeleteUserAnalytics(gomock.Any(), "test-user-1").
					Return(nil)
			},
			userID: "test-user-1",
			want:   want{err: nil},
		},
		{
			name:         "empty user ID",
			expectations: func() {},
			userID:       "",
			want:         want{err: errors.New("user ID is required")},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			err := service.DeleteUserAnalytics(ctx, tc.userID)

			assert.Equal(t, tc.want.err, err)
		})
	}
}

func TestProcessEvent(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := NewMockRepository(ctrl)
	service := NewAnalyticsService(repoMock)

	now := time.Now()
	validEvent := &Event{
		ID:        "event-1",
		EventType: "tweet_created",
		Handler:   "user-1",
		Timestamp: now,
	}

	tt := []struct {
		name         string
		expectations func()
		event        *Event
		want         error
	}{
		{
			name: "success",
			expectations: func() {
				repoMock.EXPECT().
					ProcessEvent(gomock.Any(), validEvent).
					Return(nil)
			},
			event: validEvent,
			want:  nil,
		},
		{
			name:         "nil event",
			expectations: func() {},
			event:        nil,
			want:         errors.New("event cannot be nil"),
		},
		{
			name:         "missing user ID",
			expectations: func() {},
			event: &Event{
				ID:        "event-1",
				EventType: "tweet_created",
				Timestamp: now,
			},
			want: errors.New("user ID is required in event"),
		},
		{
			name:         "missing event type",
			expectations: func() {},
			event: &Event{
				ID:        "event-1",
				Handler:   "user-1",
				Timestamp: now,
			},
			want: errors.New("event type is required"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			err := service.ProcessEvent(ctx, tc.event)

			assert.Equal(t, tc.want, err)
		})
	}
}

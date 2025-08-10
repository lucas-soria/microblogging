package tweets

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTweetService_CreateTweet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := NewMockRepository(ctrl)
	service := NewService(mockRepo)

	type args struct {
		req *Tweet
	}

	type want struct {
		tweet *Tweet
		err   error
	}

	tt := []struct {
		name         string
		args         args
		expectations func()
		want         want
		wantErr      bool
	}{
		{
			name: "successful tweet creation",
			args: args{
				req: &Tweet{
					Handler: "testuser",
					Content: Content{Text: "Hello, world!"},
				},
			},
			expectations: func() {
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, tweet *Tweet) (*Tweet, error) {
						assert.NotEmpty(t, tweet.ID)
						assert.Equal(t, "testuser", tweet.Handler)
						assert.Equal(t, "Hello, world!", tweet.Content.Text)
						assert.False(t, tweet.CreatedAt.IsZero())
						return tweet, nil
					}).
					Times(1)
			},
			want: want{
				tweet: &Tweet{
					Handler: "testuser",
					Content: Content{Text: "Hello, world!"},
				},
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "empty content",
			args: args{
				req: &Tweet{
					Handler: "testuser",
					Content: Content{Text: ""},
				},
			},
			expectations: func() {},
			want: want{
				tweet: nil,
				err:   errors.New("tweet content cannot be empty"),
			},
			wantErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			tweet, err := service.CreateTweet(ctx, tc.args.req)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, err, tc.want.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want.tweet.Handler, tweet.Handler)
				assert.Equal(t, tc.want.tweet.Content, tweet.Content)
			}
		})
	}
}

func TestTweetService_GetTweet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := NewMockRepository(ctrl)
	service := NewService(mockRepo)

	type want struct {
		tweet *Tweet
		err   error
	}

	tt := []struct {
		name         string
		tweetID      string
		expectations func()
		want         want
	}{
		{
			name:    "tweet found",
			tweetID: "123",
			expectations: func() {
				mockRepo.EXPECT().
					GetByID(ctx, "123").
					Return(&Tweet{
						ID:        "123",
						Handler:   "testuser",
						Content:   Content{Text: "Hello"},
						CreatedAt: mockTime(),
					}, nil).
					Times(1)
			},
			want: want{
				tweet: &Tweet{
					ID:        "123",
					Handler:   "testuser",
					Content:   Content{Text: "Hello"},
					CreatedAt: mockTime(),
				},
				err: nil,
			},
		},
		{
			name:         "empty tweet id",
			tweetID:      "",
			expectations: func() {},
			want: want{
				tweet: nil,
				err:   errors.New("tweet ID cannot be empty"),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			tweet, err := service.GetTweet(ctx, tc.tweetID)

			if tc.want.err != nil {
				assert.Error(t, err)
				assert.Equal(t, err, tc.want.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want.tweet.Handler, tweet.Handler)
				assert.Equal(t, tc.want.tweet.Content, tweet.Content)
			}
		})
	}
}

func TestTweetService_GetUserTweets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := NewMockRepository(ctrl)
	service := NewService(mockRepo)

	type want struct {
		tweets []*Tweet
		err    error
	}

	tt := []struct {
		name         string
		userID       string
		expectations func()
		want         want
	}{
		{
			name:   "user has tweets",
			userID: "user1",
			expectations: func() {
				mockRepo.EXPECT().
					GetByUserID(ctx, "user1").
					Return([]*Tweet{{
						ID:        "1",
						Handler:   "user1",
						Content:   Content{Text: "Tweet 1"},
						CreatedAt: mockTime(),
					}}, nil).
					Times(1)
			},
			want: want{
				tweets: []*Tweet{{
					ID:        "1",
					Handler:   "user1",
					Content:   Content{Text: "Tweet 1"},
					CreatedAt: mockTime(),
				}},
				err: nil,
			},
		},
		{
			name:         "empty user id",
			userID:       "",
			expectations: func() {},
			want: want{
				tweets: nil,
				err:    errors.New("user ID cannot be empty"),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			tweets, err := service.GetUserTweets(ctx, tc.userID)

			assert.Equal(t, tc.want.err, err)
			assert.Equal(t, tc.want.tweets, tweets)
		})
	}
}

func TestTweetService_DeleteTweet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := NewMockRepository(ctrl)
	service := NewService(mockRepo)

	type want struct {
		err error
	}

	tt := []struct {
		name         string
		tweetID      string
		expectations func()
		want         want
	}{
		{
			name:    "successful deletion",
			tweetID: "123",
			expectations: func() {
				mockRepo.EXPECT().
					GetByID(ctx, "123").
					Return(&Tweet{ID: "123", Handler: "testuser"}, nil).
					Times(1)
				mockRepo.EXPECT().
					Delete(ctx, "123").
					Return(nil).
					Times(1)
			},
			want: want{
				err: nil,
			},
		},
		{
			name:    "tweet not found",
			tweetID: "nonexistent",
			expectations: func() {
				mockRepo.EXPECT().
					GetByID(ctx, "nonexistent").
					Return(nil, nil).
					Times(1)
			},
			want: want{
				err: errors.New("tweet not found"),
			},
		},
		{
			name:         "empty tweet id",
			tweetID:      "",
			expectations: func() {},
			want: want{
				err: errors.New("tweet ID cannot be empty"),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			err := service.DeleteTweet(ctx, tc.tweetID)

			assert.Equal(t, tc.want.err, err)
		})
	}
}

// Helper function to provide consistent timestamps in tests
func mockTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2025-01-01T00:00:00Z")
	return t
}

package feed

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFeedService_GetUserTimeline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := NewMockRepository(ctrl)
	service := NewFeedService(mockRepo)

	type args struct {
		userID string
		limit  int
		offset int
	}

	type want struct {
		timeline *TimelineResponse
		err      error
	}

	tt := []struct {
		name         string
		expectations func()
		args         args
		want         want
	}{
		{
			name: "successful timeline retrieval",
			expectations: func() {
				mockRepo.EXPECT().
					GetUserTimeline(ctx, "user1", 20, 0).
					Return([]*Tweet{
						{
							ID:      "1",
							Handler: "user1",
							Content: Content{Text: "Hello"},
						},
					}, nil).
					Times(1)
			},
			args: args{
				userID: "user1",
				limit:  20,
				offset: 0,
			},
			want: want{
				timeline: &TimelineResponse{
					Tweets: []*Tweet{
						{
							ID:      "1",
							Handler: "user1",
							Content: Content{Text: "Hello"},
						},
					},
					NextOffset: 1,
				},
				err: nil,
			},
		},
		{
			name: "empty user ID",
			expectations: func() {
				// No expectations, should fail before calling repository
			},
			args: args{
				userID: "",
				limit:  20,
				offset: 0,
			},
			want: want{
				timeline: nil,
				err:      errors.New("user ID is required"),
			},
		},
		{
			name: "repository error",
			expectations: func() {
				mockRepo.EXPECT().
					GetUserTimeline(ctx, "user1", 20, 0).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			args: args{
				userID: "user1",
				limit:  20,
				offset: 0,
			},
			want: want{
				timeline: nil,
				err:      errors.New("database error"),
			},
		},
		{
			name: "pagination with next offset",
			expectations: func() {
				tweets := make([]*Tweet, 5)
				for i := 0; i < 5; i++ {
					tweets[i] = &Tweet{
						ID:      string(rune('f' + i)), // f, g, h, i, j
						Handler: "user1",
						Content: Content{Text: string(rune('f' + i))},
					}
				}
				mockRepo.EXPECT().
					GetUserTimeline(ctx, "user1", 5, 5).
					Return(tweets, nil).
					Times(1)
			},
			args: args{
				userID: "user1",
				limit:  5,
				offset: 5,
			},
			want: want{
				timeline: &TimelineResponse{
					Tweets: []*Tweet{
						{
							ID:        "f",
							Handler:   "user1",
							Content:   Content{Text: "f"},
							CreatedAt: time.Time{},
						},
						{
							ID:        "g",
							Handler:   "user1",
							Content:   Content{Text: "g"},
							CreatedAt: time.Time{},
						},
						{
							ID:        "h",
							Handler:   "user1",
							Content:   Content{Text: "h"},
							CreatedAt: time.Time{},
						},
						{
							ID:        "i",
							Handler:   "user1",
							Content:   Content{Text: "i"},
							CreatedAt: time.Time{},
						},
						{
							ID:        "j",
							Handler:   "user1",
							Content:   Content{Text: "j"},
							CreatedAt: time.Time{},
						},
					},
					NextOffset: 5,
				},
				err: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectations != nil {
				tc.expectations()
			}

			result, err := service.GetUserTimeline(ctx, tc.args.userID, tc.args.limit, tc.args.offset)

			assert.Equal(t, tc.want.err, err)
			assert.Equal(t, tc.want.timeline, result)
		})
	}
}

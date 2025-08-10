package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lucas-soria/microblogging/cmd/users/middleware"

	"github.com/lucas-soria/microblogging/internal/tweets"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetTweet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := tweets.NewMockRepository(ctrl)
	service := tweets.NewService(mockRepo)
	handler := NewTweetHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware())
	router.GET("/v1/tweets/:id", handler.GetTweet)

	now := time.Now().UTC()
	testTweet := &tweets.Tweet{
		ID:      "test-tweet-123",
		Handler: "test-user-123",
		Content: tweets.Content{
			Text: "This is a test tweet",
		},
		CreatedAt: now,
	}

	type args struct {
		id      string
		headers map[string]string
	}

	type want struct {
		statusCode int
		response   []byte
	}

	tt := []struct {
		name         string
		args         args
		expectations func(args args)
		want         want
	}{
		{
			name: "Get tweet successfully",
			args: args{
				id: "test-tweet-123",
				headers: map[string]string{
					"X-User-Id": "test-user-123",
				},
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					GetByID(ctx, args.id).
					Return(testTweet, nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusOK,
				response:   []byte(`{"id":"test-tweet-123","handler":"test-user-123","content":{"text":"This is a test tweet"},"created_at":"` + now.Format(time.RFC3339Nano) + `"}`),
			},
		},
		{
			name: "Missing user ID header",
			args: args{
				id:      "test-tweet-123",
				headers: map[string]string{},
			},
			expectations: func(args args) {},
			want: want{
				statusCode: http.StatusUnauthorized,
				response:   []byte(`{"error":"X-User-Id header is required"}`),
			},
		},
		{
			name: "Tweet not found",
			args: args{
				id: "non-existent-tweet",
				headers: map[string]string{
					"X-User-Id": "test-user-123",
				},
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					GetByID(ctx, args.id).
					Return(nil, nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusNotFound,
				response:   []byte(`{"error":"Tweet not found"}`),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.args)

			url := fmt.Sprintf("/v1/tweets/%s", tc.args.id)
			r := httptest.NewRequest(http.MethodGet, url, nil)
			for k, v := range tc.args.headers {
				r.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, tc.want.statusCode, w.Code)
			assert.Equal(t, tc.want.response, w.Body.Bytes())
		})
	}
}

func TestGetUserTweets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := tweets.NewMockRepository(ctrl)
	service := tweets.NewService(mockRepo)
	handler := NewTweetHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware())
	router.GET("/v1/tweets/users/:id", handler.GetUserTweets)

	now := time.Now().UTC()
	testTweets := []*tweets.Tweet{
		{
			ID:      "test-tweet-1",
			Handler: "test-user-123",
			Content: tweets.Content{
				Text: "First test tweet",
			},
			CreatedAt: now,
		},
		{
			ID:      "test-tweet-2",
			Handler: "test-user-123",
			Content: tweets.Content{
				Text: "Second test tweet",
			},
			CreatedAt: now.Add(-time.Hour),
		},
	}

	type args struct {
		userID  string
		headers map[string]string
	}

	type want struct {
		statusCode int
		response   []byte
	}

	tt := []struct {
		name         string
		args         args
		expectations func(args args)
		want         want
	}{
		{
			name: "Get user tweets successfully",
			args: args{
				userID: "test-user-123",
				headers: map[string]string{
					"X-User-Id": "test-user-123",
				},
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					GetByUserID(ctx, args.userID).
					Return(testTweets, nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusOK,
				response: []byte(`[{"id":"test-tweet-1","handler":"test-user-123","content":{"text":"First test tweet"},"created_at":"` + now.Format(time.RFC3339Nano) + `"},` +
					`{"id":"test-tweet-2","handler":"test-user-123","content":{"text":"Second test tweet"},"created_at":"` + now.Add(-time.Hour).Format(time.RFC3339Nano) + `"}]`),
			},
		},
		{
			name: "Missing user ID header",
			args: args{
				userID:  "test-user-123",
				headers: map[string]string{},
			},
			expectations: func(args args) {},
			want: want{
				statusCode: http.StatusUnauthorized,
				response:   []byte(`{"error":"X-User-Id header is required"}`),
			},
		},
		{
			name: "No tweets found",
			args: args{
				userID: "nonexistent-user",
				headers: map[string]string{
					"X-User-Id": "test-user-123",
				},
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					GetByUserID(ctx, args.userID).
					Return([]*tweets.Tweet{}, nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusOK,
				response:   []byte(`[]`),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.args)

			url := fmt.Sprintf("/v1/tweets/users/%s", tc.args.userID)
			r := httptest.NewRequest(http.MethodGet, url, nil)
			for k, v := range tc.args.headers {
				r.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, tc.want.statusCode, w.Code)
			assert.Equal(t, tc.want.response, w.Body.Bytes())
		})
	}
}

func TestDeleteTweet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := tweets.NewMockRepository(ctrl)
	service := tweets.NewService(mockRepo)
	handler := NewTweetHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware())
	router.DELETE("/v1/tweets/:id", handler.DeleteTweet)

	now := time.Now().UTC()
	testTweet := &tweets.Tweet{
		ID:      "test-tweet-123",
		Handler: "test-user-123",
		Content: tweets.Content{
			Text: "Test tweet to delete",
		},
		CreatedAt: now,
	}

	type args struct {
		tweetID string
		headers map[string]string
	}

	type want struct {
		statusCode int
		response   []byte
	}

	tt := []struct {
		name         string
		args         args
		expectations func(args args)
		want         want
	}{
		{
			name: "Delete tweet successfully",
			args: args{
				tweetID: "test-tweet-123",
				headers: map[string]string{
					"X-User-Id": "test-user-123",
				},
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					GetByID(ctx, args.tweetID).
					Return(testTweet, nil).
					Times(1)

				mockRepo.EXPECT().
					GetByID(ctx, args.tweetID).
					Return(testTweet, nil).
					Times(1)

				mockRepo.EXPECT().
					Delete(ctx, args.tweetID).
					Return(nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusNoContent,
				response:   nil,
			},
		},
		{
			name: "Missing user ID header",
			args: args{
				tweetID: "test-tweet-123",
				headers: map[string]string{},
			},
			expectations: func(args args) {},
			want: want{
				statusCode: http.StatusUnauthorized,
				response:   []byte(`{"error":"X-User-Id header is required"}`),
			},
		},
		{
			name: "Tweet not found",
			args: args{
				tweetID: "non-existent-tweet",
				headers: map[string]string{
					"X-User-Id": "test-user-123",
				},
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					GetByID(ctx, args.tweetID).
					Return(nil, nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusNotFound,
				response:   []byte(`{"error":"Tweet not found"}`),
			},
		},
		{
			name: "Unauthorized to delete tweet",
			args: args{
				tweetID: "test-tweet-123",
				headers: map[string]string{
					"X-User-Id": "different-user-456",
				},
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					GetByID(ctx, args.tweetID).
					Return(testTweet, nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusForbidden,
				response:   []byte(`{"error":"You can only delete your own tweets"}`),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.args)

			url := fmt.Sprintf("/v1/tweets/%s", tc.args.tweetID)
			r := httptest.NewRequest(http.MethodDelete, url, nil)
			for k, v := range tc.args.headers {
				r.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, tc.want.statusCode, w.Code)
			assert.Equal(t, tc.want.response, w.Body.Bytes())
		})
	}
}

func TestCreateTweet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := tweets.NewMockRepository(ctrl)
	service := tweets.NewService(mockRepo)
	handler := NewTweetHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware())
	router.POST("/v1/tweets", handler.CreateTweet)

	now := time.Now().UTC()

	type args struct {
		body    []byte
		headers map[string]string
	}

	type want struct {
		statusCode int
		tweet      []byte
	}

	tt := []struct {
		name         string
		args         args
		expectations func(args args)
		want         want
	}{
		{
			name: "Created Ok",
			args: args{
				body: []byte(`{"content":{"text":"This is a test tweet"}}`),
				headers: map[string]string{
					"X-User-Id": "test-user-123",
				},
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					Create(ctx, gomock.Any()).
					Return(&tweets.Tweet{
						Content: tweets.Content{
							Text: "This is a test tweet",
						},
						Handler:   "test-user-123",
						ID:        "test-tweet-123",
						CreatedAt: now,
					}, nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusCreated,
				tweet:      []byte(`{"id":"test-tweet-123","handler":"test-user-123","content":{"text":"This is a test tweet"},"created_at":"` + now.Format(time.RFC3339Nano) + `"}`),
			},
		},
		{
			name: "Invalid request body",
			args: args{
				body: []byte(`{"content":{"text":"This is a test tweet"}`),
				headers: map[string]string{
					"X-User-Id": "test-user-123",
				},
			},
			expectations: func(args args) {},
			want: want{
				statusCode: http.StatusBadRequest,
				tweet:      []byte(`{"error":"Invalid request body"}`),
			},
		},
		{
			name: "Missing user ID header",
			args: args{
				body:    []byte(`{"content":{"text":"This is a test tweet"}}`),
				headers: map[string]string{},
			},
			expectations: func(args args) {},
			want: want{
				statusCode: http.StatusUnauthorized,
				tweet:      []byte(`{"error":"X-User-Id header is required"}`),
			},
		},
		{
			name: "Could not create tweet",
			args: args{
				body: []byte(`{"content":{}}`),
				headers: map[string]string{
					"X-User-Id": "test-user-123",
				},
			},
			expectations: func(args args) {},
			want: want{
				statusCode: http.StatusBadRequest,
				tweet:      []byte(`{"error":"Invalid request body"}`),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.args)

			url := "/v1/tweets"
			r := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(tc.args.body))
			for k, v := range tc.args.headers {
				r.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, tc.want.statusCode, w.Code)
			assert.Equal(t, tc.want.tweet, w.Body.Bytes())
		})
	}
}

package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucas-soria/microblogging/cmd/feed/middleware"

	"github.com/lucas-soria/microblogging/internal/feed"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFeedHandler_GetUserTimeline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := feed.NewMockRepository(ctrl)
	service := feed.NewService(mockRepo)
	handler := NewFeedHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware())
	router.GET("/v1/feed/timeline", handler.GetUserTimeline)

	type args struct {
		userID      string
		queryParams string
	}

	type want struct {
		statusCode int
		response   []byte
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
					GetUserTimeline(gomock.Any(), "user1", 20, 0).
					Return([]*feed.Tweet{
						{
							ID:      "1",
							Handler: "user1",
							Content: feed.Content{Text: "Hello"},
						},
					}, nil).
					Times(1)
			},
			args: args{
				userID:      "user1",
				queryParams: "",
			},
			want: want{
				statusCode: http.StatusOK,
				response:   []byte(`{"tweets":[{"id":"1","handler":"user1","content":{"text":"Hello"},"created_at":"0001-01-01T00:00:00Z"}],"next_offset":1}`),
			},
		},
		{
			name: "with pagination",
			expectations: func() {
				mockRepo.EXPECT().
					GetUserTimeline(gomock.Any(), "user1", 10, 5).
					Return([]*feed.Tweet{}, nil).
					Times(1)
			},
			args: args{
				userID:      "user1",
				queryParams: "?limit=10&offset=5",
			},
			want: want{
				statusCode: http.StatusOK,
				response:   []byte(`{"tweets":[],"next_offset":0}`),
			},
		},
		{
			name:         "missing user id header",
			expectations: func() {},
			args: args{
				userID:      "",
				queryParams: "",
			},
			want: want{
				statusCode: http.StatusUnauthorized,
				response:   []byte(`{"error":"X-User-Id header is required"}`),
			},
		},
		{
			name: "service error",
			expectations: func() {
				mockRepo.EXPECT().
					GetUserTimeline(gomock.Any(), "user1", 20, 0).
					Return(nil, assert.AnError).
					Times(1)
			},
			args: args{
				userID:      "user1",
				queryParams: "",
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				response:   []byte(`{"error":"Failed to get user timeline"}`),
			},
		},
		{
			name: "invalid limit",
			args: args{
				userID:      "user1",
				queryParams: "?limit=invalid",
			},
			expectations: func() {},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   []byte(`{"error":"Invalid limit parameter"}`),
			},
		},
		{
			name: "invalid offset",
			args: args{
				userID:      "user1",
				queryParams: "?offset=invalid",
			},
			expectations: func() {},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   []byte(`{"error":"Invalid offset parameter"}`),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			url := fmt.Sprintf("/v1/feed/timeline%s", tc.args.queryParams)
			r := httptest.NewRequest(http.MethodGet, url, nil)
			r.Header.Set("X-User-Id", tc.args.userID)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, tc.want.statusCode, w.Code)
			assert.Equal(t, string(tc.want.response), w.Body.String())
		})
	}
}

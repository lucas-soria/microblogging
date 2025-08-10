package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucas-soria/microblogging/cmd/users/middleware"

	"github.com/lucas-soria/microblogging/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := users.NewMockRepository(ctrl)
	service := users.NewService(mockRepo)
	handler := NewUserHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/v1/users", handler.CreateUser)

	type args struct {
		body []byte
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
			name: "Create user successfully",
			args: args{
				body: []byte(`{"handler":"testuser","first_name":"Test","last_name":"User"}`),
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					CreateUser(ctx, &users.User{
						Handler:   "testuser",
						FirstName: "Test",
						LastName:  "User",
					}).
					Return(nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusCreated,
				response:   []byte(`{"handler":"testuser","first_name":"Test","last_name":"User"}`),
			},
		},
		{
			name: "Invalid request body",
			args: args{
				body: []byte(`{"handler":"","first_name":"Test",}`),
			},
			expectations: func(args args) {},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   []byte(`{"error":"Invalid request body"}`),
			},
		},
		{
			name: "User already exists",
			args: args{
				body: []byte(`{"handler":"existinguser","first_name":"Test","last_name":"User"}`),
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					CreateUser(ctx, gomock.Any()).
					Return(users.ErrHandlerExists).
					Times(1)
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				response:   []byte(`{"error":"failed to create user"}`),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.args)

			r := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(tc.args.body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.want.statusCode, w.Code)
			assert.Equal(t, tc.want.response, w.Body.Bytes())
		})
	}
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := users.NewMockRepository(ctrl)
	service := users.NewService(mockRepo)
	handler := NewUserHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/v1/users/:id", handler.GetUser)

	type args struct {
		userID string
	}

	type want struct {
		statusCode int
		response   []byte
	}

	testUser := &users.User{
		Handler:   "testuser",
		FirstName: "Test",
		LastName:  "User",
	}

	tt := []struct {
		name         string
		args         args
		expectations func(args args)
		want         want
	}{
		{
			name: "Get user successfully",
			args: args{
				userID: "testuser",
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					GetUser(ctx, args.userID).
					Return(testUser, nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusOK,
				response:   []byte(`{"handler":"testuser","first_name":"Test","last_name":"User"}`),
			},
		},
		{
			name: "User not found",
			args: args{
				userID: "nonexistent",
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					GetUser(ctx, args.userID).
					Return(nil, users.ErrUserNotFound).
					Times(1)
			},
			want: want{
				statusCode: http.StatusNotFound,
				response:   []byte(`{"error":"user not found"}`),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.args)

			url := fmt.Sprintf("/v1/users/%s", tc.args.userID)
			r := httptest.NewRequest(http.MethodGet, url, nil)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.want.statusCode, w.Code)
			assert.Equal(t, tc.want.response, w.Body.Bytes())
		})
	}
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := users.NewMockRepository(ctrl)
	service := users.NewService(mockRepo)
	handler := NewUserHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware())
	router.DELETE("/v1/users/:id", handler.DeleteUser)

	type args struct {
		userID string
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
			name: "Delete user successfully",
			args: args{
				userID: "testuser",
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					DeleteUser(ctx, args.userID).
					Return(nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusNoContent,
				response:   nil,
			},
		},
		{
			name: "User not found",
			args: args{
				userID: "nonexistent",
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					DeleteUser(ctx, args.userID).
					Return(users.ErrUserNotFound).
					Times(1)
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				response:   []byte(`{"error":"failed to delete user"}`),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.args)

			url := fmt.Sprintf("/v1/users/%s", tc.args.userID)
			r := httptest.NewRequest(http.MethodDelete, url, nil)
			r.Header.Set("X-User-Id", tc.args.userID)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.want.statusCode, w.Code)
			assert.Equal(t, tc.want.response, w.Body.Bytes())
		})
	}
}

func TestFollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := users.NewMockRepository(ctrl)
	service := users.NewService(mockRepo)
	handler := NewUserHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware())
	router.POST("/v1/users/:id/follow", handler.FollowUser)

	type args struct {
		followeeID string
		headers    map[string]string
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
			name: "Follow user successfully",
			args: args{
				followeeID: "user2",
				headers: map[string]string{
					"X-User-Id": "user1",
				},
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					FollowUser(ctx, "user1", "user2").
					Return(nil).
					Times(1)
			},
			want: want{
				statusCode: http.StatusAccepted,
				response:   nil,
			},
		},
		{
			name: "Missing user ID header",
			args: args{
				followeeID: "user2",
				headers:    map[string]string{},
			},
			expectations: func(args args) {},
			want: want{
				statusCode: http.StatusUnauthorized,
				response:   []byte(`{"error":"X-User-Id header is required"}`),
			},
		},
		{
			name: "User not found",
			args: args{
				followeeID: "nonexistent",
				headers: map[string]string{
					"X-User-Id": "user1",
				},
			},
			expectations: func(args args) {
				mockRepo.EXPECT().
					FollowUser(ctx, "user1", "nonexistent").
					Return(users.ErrUserNotFound).
					Times(1)
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				response:   []byte(`{"error":"failed to follow user"}`),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations(tc.args)

			url := fmt.Sprintf("/v1/users/%s/follow", tc.args.followeeID)
			r := httptest.NewRequest(http.MethodPost, url, nil)
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

// Similar tests for UnfollowUser, GetUserFollowers, and GetUserFollowing would follow the same pattern
// as TestFollowUser, testing various scenarios like success, missing auth, and not found cases.

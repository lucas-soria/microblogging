package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucas-soria/microblogging/internal/analytics"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserAnalytics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := analytics.NewMockRepository(ctrl)
	service := analytics.NewAnalyticsService(repoMock)
	handler := NewAnalyticsHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/v1/analytics/users/:id", handler.GetUserAnalytics)

	userID := uuid.New().String()

	type want struct {
		statusCode int
		response   []byte
	}

	tt := []struct {
		name         string
		userID       string
		expectations func()
		want         want
	}{
		{
			name:   "success",
			userID: userID,
			expectations: func() {
				repoMock.EXPECT().
					GetUserAnalytics(gomock.Any(), userID).
					Return(&analytics.UserAnalytics{
						Handler:      userID,
						IsInfluencer: true,
						IsActive:     true,
					}, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				response:   []byte(`{"handler":"` + userID + `","is_influencer":true,"is_active":true}`),
			},
		},
		{
			name:   "not found",
			userID: userID,
			expectations: func() {
				repoMock.EXPECT().
					GetUserAnalytics(gomock.Any(), userID).
					Return(nil, errors.New("user analytics not found"))
			},
			want: want{
				statusCode: http.StatusNotFound,
				response:   []byte(`{"error":"user analytics not found"}`),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			req, err := http.NewRequest(http.MethodGet, "/v1/analytics/users/"+tc.userID, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.want.statusCode, rr.Code)
			assert.Equal(t, string(tc.want.response), rr.Body.String())
		})
	}
}

func TestGetAllUserAnalytics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := analytics.NewMockRepository(ctrl)
	service := analytics.NewAnalyticsService(repoMock)
	handler := NewAnalyticsHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/v1/analytics/users", handler.GetAllUserAnalytics)

	userID1 := uuid.New().String()
	userID2 := uuid.New().String()

	type want struct {
		statusCode int
		response   []byte
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
					Return([]*analytics.UserAnalytics{
						{
							Handler:      userID1,
							IsInfluencer: true,
							IsActive:     true,
						},
						{
							Handler:      userID2,
							IsInfluencer: false,
							IsActive:     true,
						},
					}, nil)
			},
			want: want{
				statusCode: http.StatusOK,
				response:   []byte(`[{"handler":"` + userID1 + `","is_influencer":true,"is_active":true},{"handler":"` + userID2 + `","is_influencer":false,"is_active":true}]`),
			},
		},
		{
			name: "internal server error",
			expectations: func() {
				repoMock.EXPECT().
					GetAllUserAnalytics(gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				response:   []byte(`{"error":"failed to get user analytics"}`),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			req, err := http.NewRequest(http.MethodGet, "/v1/analytics/users", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.want.statusCode, rr.Code)
			assert.Equal(t, string(tc.want.response), rr.Body.String())
		})
	}
}

func TestDeleteUserAnalytics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := analytics.NewMockRepository(ctrl)
	service := analytics.NewAnalyticsService(repoMock)
	handler := NewAnalyticsHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/v1/analytics/users/:id", handler.DeleteUserAnalytics)

	userID := uuid.New().String()

	type want struct {
		statusCode int
		response   []byte
	}

	tt := []struct {
		name         string
		userID       string
		expectations func()
		want         want
	}{
		{
			name:   "success",
			userID: userID,
			expectations: func() {
				repoMock.EXPECT().
					DeleteUserAnalytics(gomock.Any(), userID).
					Return(nil)
			},
			want: want{
				statusCode: http.StatusNoContent,
				response:   []byte(``),
			},
		},
		{
			name:   "not found",
			userID: userID,
			expectations: func() {
				repoMock.EXPECT().
					DeleteUserAnalytics(gomock.Any(), userID).
					Return(errors.New("user not found"))
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				response:   []byte(`{"error":"user not found"}`),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			req, err := http.NewRequest(http.MethodDelete, "/v1/analytics/users/"+tc.userID, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.want.statusCode, rr.Code)
			assert.Equal(t, string(tc.want.response), rr.Body.String())
		})
	}
}

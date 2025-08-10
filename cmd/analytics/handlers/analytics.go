package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-soria/microblogging/internal/analytics"
)

// AnalyticsHandler handles HTTP requests for analytics
type AnalyticsHandler struct {
	service analytics.Service
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(service analytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

// GetUserAnalyticsResponse represents the response for GetUserAnalytics
type GetUserAnalyticsResponse struct {
	Handler      string `json:"handler"`
	IsInfluencer bool   `json:"is_influencer"`
	IsActive     bool   `json:"is_active"`
}

// GetUserAnalytics handles GET /v1/analytics/users/:id
func (handler *AnalyticsHandler) GetUserAnalytics(ctx *gin.Context) {
	userID := ctx.Param("id")

	// Call service
	analytics, err := handler.service.GetUserAnalytics(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return response
	response := GetUserAnalyticsResponse{
		Handler:      analytics.Handler,
		IsInfluencer: analytics.IsInfluencer,
		IsActive:     analytics.IsActive,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetAllUserAnalytics handles GET /v1/analytics/users
func (handler *AnalyticsHandler) GetAllUserAnalytics(c *gin.Context) {
	// Call service
	analyticsList, err := handler.service.GetAllUserAnalytics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user analytics"})
		return
	}

	// Convert to response type
	response := make([]GetUserAnalyticsResponse, 0, len(analyticsList))
	for _, a := range analyticsList {
		response = append(response, GetUserAnalyticsResponse{
			Handler:      a.Handler,
			IsInfluencer: a.IsInfluencer,
			IsActive:     a.IsActive,
		})
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUserAnalytics handles DELETE /v1/analytics/users/:id
func (handler *AnalyticsHandler) DeleteUserAnalytics(ctx *gin.Context) {
	userID := ctx.Param("id")

	// Call service
	err := handler.service.DeleteUserAnalytics(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

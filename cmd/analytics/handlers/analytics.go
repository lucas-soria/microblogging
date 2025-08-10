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
func (h *AnalyticsHandler) GetUserAnalytics(c *gin.Context) {
	userID := c.Param("id")

	// Call service
	analytics, err := h.service.GetUserAnalytics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return response
	response := GetUserAnalyticsResponse{
		Handler:      analytics.Handler,
		IsInfluencer: analytics.IsInfluencer,
		IsActive:     analytics.IsActive,
	}

	c.JSON(http.StatusOK, response)
}

// GetAllUserAnalytics handles GET /v1/analytics/users
func (h *AnalyticsHandler) GetAllUserAnalytics(c *gin.Context) {
	// Call service
	analyticsList, err := h.service.GetAllUserAnalytics(c.Request.Context())
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
func (h *AnalyticsHandler) DeleteUserAnalytics(c *gin.Context) {
	userID := c.Param("id")

	// Call service
	err := h.service.DeleteUserAnalytics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
